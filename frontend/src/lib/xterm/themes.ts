import { colorScheme } from "../util";
import { type ITheme } from "@xterm/xterm";

export const lightTheme: ITheme = {
  background: colorScheme.light["base-100"],
  foreground: colorScheme.light["base-content"],
  cursor: colorScheme.light.primary,
  cursorAccent: colorScheme.light["base-300"],
  selectionBackground: colorScheme.light.primary,
  selectionForeground: colorScheme.light["base-content"],
};

export const darkTheme: ITheme = {
  background: colorScheme.dark["base-100"],
  foreground: colorScheme.dark["base-content"],
  cursor: colorScheme.dark.primary,
  cursorAccent: colorScheme.dark["base-300"],
  selectionBackground: colorScheme.dark.primary,
  selectionForeground: colorScheme.dark["base-content"],
};
