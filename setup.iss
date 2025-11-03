; ==================================================================
; Inno Setup Script for TAMALabs (Normal App + Auto-Start Tray)
; ==================================================================

; Read version from version.go
#define AppVer ReadIni(SourcePath + "version.ini", "Version", "AppVersion", "v1.0.0")

[Setup]
AppId={{A8A93F44-8D2B-4D75-9CC8-0C52B2184AC4}}
AppName=TAMALabs
AppVersion={#AppVer}
AppPublisher=Elga Tama
AppPublisherURL=https://tamalabs.biosystems.id/
AppSupportURL=https://tamalabs.biosystems.id/support
AppUpdatesURL=https://tamalabs.biosystems.id/updates
DefaultDirName={autopf}\TAMALabs
DefaultGroupName=TAMALabs
AllowNoIcons=yes
OutputBaseFilename=TAMALabs-setup-{#AppVer}
OutputDir=.\installers
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Dirs]
Name: "{app}\logs"; Permissions: users-modify

[Files]
Source: "bin\TAMALabs.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: ".env"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\TAMALabsTray.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
; Optional desktop shortcut
Name: "{autodesktop}\TAMALabs (Open Web)"; \
    Filename: "rundll32.exe"; \
    Parameters: "url.dll,FileProtocolHandler http://127.0.0.1:8322"; \
    IconFilename: "{app}\TAMALabs.exe"; \
    Tasks: desktopicon

; Auto-start main app via Startup folder (for all users)
Name: "{commonstartup}\TAMALabs"; \
    Filename: "{app}\TAMALabs.exe"; \
    WorkingDir: "{app}"; \
    IconFilename: "{app}\TAMALabs.exe"

; Auto-start tray via Startup folder (for all users)
Name: "{commonstartup}\TAMALabs Tray"; \
    Filename: "{app}\TAMALabsTray.exe"; \
    WorkingDir: "{app}"; \
    IconFilename: "{app}\TAMALabsTray.exe"

[Run]
; Jalankan aplikasi utama langsung (tanpa muncul checkbox di akhir)
Filename: "{app}\TAMALabs.exe"; \
    Description: "Start TAMALabs"; \
    Flags: nowait skipifsilent; \
    StatusMsg: "Starting TAMALabs..."

; Jalankan tray langsung (tanpa muncul checkbox di akhir)
Filename: "{app}\TAMALabsTray.exe"; \
    Description: "Start TAMALabs Tray"; \
    Flags: nowait skipifsilent; \
    WorkingDir: "{app}"; \
    StatusMsg: "Starting TAMALabs Tray..."

[UninstallRun]
; Tutup tray dan app sebelum uninstall
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabs.exe /T"; Flags: runhidden waituntilterminated; RunOnceId: "KillMain"
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabsTray.exe /T"; Flags: runhidden waituntilterminated; RunOnceId: "KillTray"

; Hapus folder aplikasi
Filename: "cmd.exe"; Parameters: "/C rmdir /S /Q ""{app}"""; Flags: runhidden waituntilterminated; RunOnceId: "RemoveAppFolder"

[UninstallDelete]
; Hapus shortcut startup & desktop
Type: files; Name: "{commonstartup}\TAMALabs Tray.lnk"
Type: files; Name: "{autodesktop}\TAMALabs (Open Web).lnk"
Type: dirifempty; Name: "{app}\logs"
