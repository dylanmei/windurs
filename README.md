windurs
=======

shitty cli for interacting with windows machines via winrm.

Example:

    make
    bin/windurs --help
    bin/windurs ls -addr=mywindows:5985 -user=vagrant -pass=vagrant C:/Windows
    bin/windurs cmd -addr=mywindows:5985 -user=vagrant -pass=vagrant echo %COMPUTERNAME%
    bin/windurs info -addr=mywindows:5985 -user=vagrant -pass=vagrant info
