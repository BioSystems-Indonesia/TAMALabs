; -- Inno Setup Script for TAMALabs with Windows Service Registration --
; ==============================================
; FEATURES:
;  - Automatic service installation and startup
;  - System tray auto-start with administrator privileges (via scheduled task)
;  - No post-install dialog prompts
;  - Clean uninstallation with service and task removal
; ==============================================

[Setup]
AppId={{F4A4A2A2-702D-4B1F-A88E-5E3A1A8E2E8A}}
AppName=TAMALabs
AppVersion=1.0
AppPublisher=Elga Tama
AppPublisherURL=https://www.elgatama.com/
AppSupportURL=https://www.elgatama.com/support
AppUpdatesURL=https://www.elgatama.com/updates
DefaultDirName={autopf}\TAMALabs
DefaultGroupName=TAMALabs
AllowNoIcons=yes
OutputBaseFilename=TAMALabs-setup-v1.0
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
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Dirs]
Name: "{app}\tmp"; Permissions: users-modify

[Files]
; Main app binary
Source: "bin\TAMALabs.exe"; DestDir: "{app}"; Flags: ignoreversion
; Include NSSM for service management
Source: "bin\nssm.exe"; DestDir: "{app}"; Flags: ignoreversion
; Include Tray for TAMALabs tray system
Source: "bin\TAMALabsTray.exe"; DestDir: "{app}"; Flags: ignoreversion
; Include Service Helper with admin privileges
Source: "bin\service-helper.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
; Shortcut yang membuka browser default ke http://127.0.0.1:8322
Name: "{autodesktop}\TAMALabs (Open Web)"; \
    Filename: "rundll32.exe"; \
    Parameters: "url.dll,FileProtocolHandler http://127.0.0.1:8322"; \
    IconFilename: "{app}\TAMALabs.exe"; \
    Tasks: desktopicon

; Backup startup shortcut (fallback jika scheduled task gagal)
Name: "{userstartup}\TAMALabs Tray"; \
    Filename: "{app}\TAMALabsTray.exe"; \
    WorkingDir: "{app}"

[Run]
; --- Register TAMALabs as Windows Service ---
; Install the service (if not already installed)
Filename: "{app}\nssm.exe"; \
    Parameters: "install TAMALabs ""{app}\TAMALabs.exe"""; \
    Flags: runhidden waituntilterminated; \
    StatusMsg: "Registering TAMALabs service..."

; Set the service startup type to automatic
Filename: "sc.exe"; \
    Parameters: "config TAMALabs start= auto"; \
    Flags: runhidden waituntilterminated; \
    StatusMsg: "Configuring TAMALabs service startup..."

; Start the service immediately
Filename: "{app}\nssm.exe"; \
    Parameters: "start TAMALabs"; \
    Flags: runhidden waituntilterminated; \
    StatusMsg: "Starting TAMALabs service..."

; Create scheduled task for auto-start tray
Filename: "schtasks.exe"; \
    Parameters: "/create /tn ""TAMALabs Tray"" /tr ""{app}\TAMALabsTray.exe"" /sc onlogon /f"; \
    Flags: runhidden waituntilterminated; \
    StatusMsg: "Setting up TAMALabs tray auto-start..."

; Start TAMALabsTray immediately (now safe with asInvoker manifest)
Filename: "{app}\TAMALabsTray.exe"; \
    Flags: runasoriginaluser nowait; \
    StatusMsg: "Starting TAMALabs system tray..."

; Start TAMALabsTray immediately as current user (now safe with asInvoker manifest)
Filename: "{app}\TAMALabsTray.exe"; \
    Flags: runasoriginaluser nowait; \
    StatusMsg: "Starting TAMALabs system tray..."

[UninstallRun]
; --- Stop and remove service when uninstalling ---
Filename: "{app}\nssm.exe"; Parameters: "stop TAMALabs"; Flags: runhidden waituntilterminated; RunOnceId: "StopService"
Filename: "{app}\nssm.exe"; Parameters: "remove TAMALabs confirm"; Flags: runhidden waituntilterminated; RunOnceId: "RemoveService"
; Remove scheduled task for tray auto-start (ignore errors if not exists)
Filename: "schtasks.exe"; Parameters: "/delete /tn ""TAMALabs Tray"" /f"; Flags: runhidden; RunOnceId: "RemoveScheduledTask"

[UninstallDelete]
; Clean up any startup shortcuts
Type: files; Name: "{userstartup}\TAMALabs Tray.lnk"
