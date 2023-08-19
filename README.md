# Camunda Incident Agregator

Small little golang app for aggregating incidents (unhandled exceptions) of one or multiple
[Camunda][camunda] instances.

---

## Build With

![golang][goBadge]

## Build for

![Windows][windows-Badge]  ![Linux][linux-Badge]  ![MacOS][macOS-Badge]

---

## What does this app do ?

This application aggregates the incidents of one or multiple [Camunda][camunda] instances  
depending on a given time window or frame.

### What time frames are supported ?

in short u can search for:

- Day
- Week
- Month



---

## Configuration

This Application supports a environment/config file.  
The file is located ad resources/config/environment.yaml.  
In this case I decided to use a .yaml file.

The configuration should load automatically from

- ./resources/config/
- ../resources/config/
- .

depending on the run directory of the build program.

> **IN CASE OF NO AUTO CONFIG DETECTION**
>
>you can use the -dir flag to give a custom directory for the config

---

## Logs

Logs are found after starting in a new directory ./Logs

---

## Export

The application can also export a csv file for each configured [camunda][camunda] instance.  
The exports can be found in the ./Export directory.

---

## Build

To manually build u can use the normal go build command like  
`go build -C cmd/main -o Path/to/youre/build/application_name`  
or u can run the `build.sh` shell script for building the application multi-platform
> NOTICE  
> that the shell script does not build all possible platforms.  
> Currently only Linux,Windows and macOS are supported in there  
> x86/64 and ARM64 flavors

[camunda]: https://camunda.com/

[goBadge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white

[linux-Badge]: https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black

[windows-Badge]: https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white

[macOS-Badge]:    https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white
