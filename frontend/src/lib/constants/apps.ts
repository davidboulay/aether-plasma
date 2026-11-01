type SpecialAppFlag = 'includeZed' | 'includeVscode' | 'includeNeovim';

// These legacy flags still drive editor-specific install steps. Keep them in
// sync with the explicit app inclusion set.
export const SPECIAL_APP_FLAGS: Record<string, SpecialAppFlag> = {
    zed: 'includeZed',
    vscode: 'includeVscode',
    neovim: 'includeNeovim',
};

export const SPECIAL_APP_KEYS = new Set(Object.keys(SPECIAL_APP_FLAGS));

// Order in which specials should appear in UI lists.
export const SPECIAL_APP_ORDER: readonly string[] = ['neovim', 'zed', 'vscode'];

// Templates that are always written and never user-togglable. Hidden from
// the Apps sidebar and the Targets strip. colors.toml is the canonical
// machine-readable export consumed by external tools, so opting out of it
// would silently break integrations.
export const ALWAYS_INCLUDED_APPS = new Set(['colors', 'kde']);

// Aether Plasma: targets for the Hyprland/Omarchy stack, which a KDE Plasma
// session doesn't run. Hidden from the Apps sidebar and the Targets strip.
export const PLASMA_HIDDEN_APPS = new Set([
    'hyprland',
    'hyprlock',
    'mako',
    'swayosd',
    'walker',
    'waybar',
    'wofi',
]);

// Aether Plasma integrations that are on unless explicitly turned off.
// Keep in sync with plasmaDefaultApps in internal/theme/writer.go.
export const DEFAULT_ON_APPS = new Set([
    'icons',
    'konsole',
    'ghostty',
    'alacritty',
    'kitty',
    'foot',
    'btop',
]);

// Everything the target lists leave out.
export function isHiddenApp(key: string): boolean {
    return ALWAYS_INCLUDED_APPS.has(key) || PLASMA_HIDDEN_APPS.has(key);
}

// Pretty labels for known templates. Anything not listed falls back to a
// title-cased version of the key.
const APP_LABELS: Record<string, string> = {
    zed: 'Zed',
    vscode: 'VS Code',
    neovim: 'Neovim',
    alacritty: 'Alacritty',
    btop: 'Btop',
    chromium: 'Chromium',
    colors: 'Colors (.toml)',
    foot: 'Foot',
    ghostty: 'Ghostty',
    hyprland: 'Hyprland',
    hyprlock: 'Hyprlock',
    icons: 'Icons',
    kde: 'KDE Plasma',
    kitty: 'Kitty',
    konsole: 'Konsole',
    mako: 'Mako',
    swayosd: 'SwayOSD',
    vencord: 'Vencord',
    walker: 'Walker',
    warp: 'Warp',
    waybar: 'Waybar',
    wofi: 'Wofi',
    zellij: 'Zellij',
};

export function appLabel(key: string): string {
    return APP_LABELS[key] ?? key.charAt(0).toUpperCase() + key.slice(1);
}
