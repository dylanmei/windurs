package windurs

import "text/template"

type elevatedOptions struct {
	User           string
	Password       string
	LogFile        string
	Description    string
	EncodedCommand string
}

var elevatedTemplate = template.Must(template.New("ElevatedCommand").Parse(`
$task_name = "packer-elevated-shell"
$log_file = "$env:TEMP\{{.LogFile}}"
$task_xml = @'
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
	<Description>{{.Description}}</Description>
  </RegistrationInfo>
  <Principals>
    <Principal id="Author">
      <UserId>{{.User}}</UserId>
      <LogonType>Password</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT2H</ExecutionTimeLimit>
    <Priority>4</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>cmd</Command>
	  <Arguments>/c powershell.exe -EncodedCommand {{.EncodedCommand}} &gt; %TEMP%\{{.LogFile}} 2&gt;&amp;1</Arguments>
    </Exec>
  </Actions>
</Task>
'@
$schedule = New-Object -ComObject "Schedule.Service"
$schedule.Connect()
$task = $schedule.NewTask($null)
$task.XmlText = $task_xml
$folder = $schedule.GetFolder("\")
$folder.RegisterTaskDefinition($task_name, $task, 6, "{{.User}}", "{{.Password}}", 1, $null) | Out-Null
$registered_task = $folder.GetTask("\$task_name")
$registered_task.Run($null) | Out-Null
$timeout = 10
$sec = 0
while ( (!($registered_task.state -eq 4)) -and ($sec -lt $timeout) ) {
  Start-Sleep -s 1
  $sec++
}
function SlurpOutput($log_file, $cur_line) {
  if (Test-Path $log_file) {
    Get-Content $log_file | select -skip $cur_line | ForEach {
      $cur_line += 1
      Write-Host "$_"
    }
  }
  return $cur_line
}
$cur_line = 0
do {
  Start-Sleep -m 100
  $cur_line = SlurpOutput $log_file $cur_line
} while (!($registered_task.state -eq 3))
$exit_code = $registered_task.LastTaskResult
[System.Runtime.Interopservices.Marshal]::ReleaseComObject($schedule) | Out-Null
exit $exit_code`))
