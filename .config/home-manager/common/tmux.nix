{ pkgs, clipboardCommand, ... }:
{
  programs.tmux = {
    enable = true;
    terminal = "xterm-256color";
    baseIndex = 1;
    escapeTime = 1;
    historyLimit = 999999999;
    keyMode = "vi";
    prefix = "C-q";
    secureSocket = true;
    shell = "${pkgs.zsh}/bin/zsh";
    newSession = true;
    customPaneNavigationAndResize = true;
    resizeAmount = 5;
    disableConfirmationPrompt = true;
    mouse = true;

    plugins = with pkgs; [
      tmuxPlugins.open
      tmuxPlugins.sensible
      tmuxPlugins.urlview
      tmuxPlugins.pain-control
      tmuxPlugins.jump
      tmuxPlugins.resurrect
      {
        plugin = tmuxPlugins.continuum;
        extraConfig = ''
          set -g @continuum-restore 'on'
        '';
      }
      {
        plugin = tmuxPlugins.dracula;
        extraConfig = ''
          set -g clock-mode-colour "#81a2be"
          set-option -g @dracula-plugins "battery"
        '';
      }
    ];

    extraConfig = ''
      set-option -g default-command $SHELL
      set-option -g display-panes-time 15000
      set-option -g status-position top
      set-option -g status-interval 1

      bind-key C-g display-panes

      bind -n WheelUpPane if-shell -F -t = "#{mouse_any_flag}" "send-keys -M" "if -Ft= '#{pane_in_mode}' 'send-keys -M' 'copy-mode -e'"

      # clipboard tool is provided per-host (hosts/wsl.nix, hosts/macos.nix)
      set -s copy-command "${clipboardCommand}"
      # copy-command only covers tmux's own copy-mode. Applications inside a pane
      # (nvim, TUIs, remote SSH) copy via OSC 52, which tmux's default of
      # "external" drops outright -- it neither buffers nor forwards it. "on"
      # accepts the sequence and relays it to the outer terminal.
      set -g set-clipboard on
      # some tools wrap OSC 52 in tmux's DCS passthrough instead of emitting it
      # bare; without this they are dropped even with set-clipboard on.
      set -g allow-passthrough on
      bind-key -T copy-mode-vi v send-keys -X begin-selection
      bind-key -T copy-mode-vi y send-keys -X copy-pipe
      bind-key -T copy-mode-vi Enter send-keys -X copy-pipe-and-cancel
    '';
  };
}
