; -- Inno Setup Script for LIS Elgatama --
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value for other applications.
; (To generate a new GUID, click Tools | Generate GUID in the IDE.)
AppId={{F4A4A2A2-702D-4B1F-A88E-5E3A1A8E2E8A}}
AppName=LIS Elgatama
AppVersion=1.0
ArchitecturesInstallIn64BitMode=x64compatible
;AppVerName=LIS Elgatama 1.0
AppPublisher=Elgatama
AppPublisherURL=https://www.elgatama.com/
AppSupportURL=https://www.elgatama.com/support
AppUpdatesURL=https://www.elgatama.com/updates
DefaultDirName={autopf}\LIS Elgatama
DefaultGroupName=LIS Elgatama
AllowNoIcons=yes
LicenseFile=
InfoBeforeFile=
InfoAfterFile=
; "OutputBaseFilename" is the name of the generated Setup file.
OutputBaseFilename=lis-elgatama-setup-v1.0
; The output directory for the compiled setup file.
OutputDir=.\installers
Compression=lzma
SolidCompression=yes
WizardStyle=modern


[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
; This creates a checkbox in the wizard allowing the user to create a desktop icon.
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Dirs]
; This entry creates a 'tmp' folder inside the application's installation directory.
Name: "{app}\tmp"; Permissions: users-modify

[Files]
; This is the main application binary.
; Source: The path to your compiled Go application.
; DestDir: "{app}" is the folder the user selects during installation (e.g., C:\Program Files\LIS Elgatama)
; Flags: ignoreversion tells the installer to overwrite the file regardless of version numbers.
Source: "bin\winapp.exe"; DestDir: "{app}"; Flags: ignoreversion

; NOTE: If your application requires other files or directories (e.g., config files, assets),
; you can add them here. For example, to include everything in a 'config' folder:
; Source: "config\*"; DestDir: "{app}\config"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
; Icon for the Start Menu programs list.
Name: "{group}\LIS Elgatama"; Filename: "{app}\winapp.exe"
; Optional: Icon for the Start Menu to uninstall the application.
Name: "{group}\{cm:UninstallProgram,LIS Elgatama}"; Filename: "{uninstallexe}"
; Desktop icon, created only if the "desktopicon" task was checked.
Name: "{autodesktop}\LIS Elgatama"; Filename: "{app}\winapp.exe"; Tasks: desktopicon

[Run]
; This gives the user an option to run the application right after installation is complete.
Filename: "{app}\winapp.exe"; Description: "{cm:LaunchProgram,LIS Elgatama}"; Flags: nowait postinstall skipifsilent
