set curDir=%~dp0
set output="%curDir%build\tank_{{.OS}}_{{.Arch}}"
cd ./src && gox -os="windows darwin" -output %output%