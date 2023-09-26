import { EditorView } from "@codemirror/view";
import type { Extension } from "@codemirror/state";
import { HighlightStyle, syntaxHighlighting } from "@codemirror/language";
import { tags as t } from "@lezer/highlight";
import { yamlLanguage } from "@codemirror/lang-yaml";

const chalky = "#e5c07b",
  coral = "#e06c75",
  cyan = "#56b6c2",
  stone = "#7d8799",
  malibu = "#61afef",
  sage = "#98c379",
  whiskey = "#d19a66",
  violet = "#c678dd";

export const cmTheme: Extension = EditorView.theme({
  "&.cm-editor": {
    backgroundColor: "oklch(var(--b1))",
    "font-size": "15px",
    height: "100%",
    width: "100%",
    overflow: "hidden",
    "border-radius": "6px",
  },

  ".cm-gutter": {
    height: "100%",
    color: "oklch(var(--bc) / 0.4)",
    border: "none",
  },

  ".cm-gutters": {
    "border-style": "none",
    "justify-content": "flex-end",
    "background-color": "oklch(var(--b2))",
  },

  "&.cm-editor .cm-activeLineGutter": {
    "background-color": "oklch(var(--b2) / 0.4)",
    color: "oklch(var(--bc))",
  },

  "&.cm-editor .cm-activeLine": {
    "background-color": "oklch(var(--b2) / 0.4)",
  },

  ".cm-scroller": {
    height: "100%",
    overflow: "auto",
  },

  ".cm-cursorLayer .cm-cursor": {
    "border-left": "6px solid oklch(var(--p))",
  },

  "&.cm-editor.cm-focused": {
    outline: "none",
  },

  "&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection":
    {
      "background-color": "oklch(var(--p) / 0.6)",
    },

  ".cm-foldPlaceholder": {
    color: "oklch(var(--bc))",
    "background-color": "oklch(var(--b2))",
    border: "solid",
    "border-width": "1px",
    "border-color": "oklch(var(--bc) / 0.6)",
    "margin-left": "0.5rem",
    "margin-right": "0.5rem",
  },

  ".cm-foldGutter > .cm-gutterElement": {
    color: "oklch(var(--a))",
    "margin-right": "0.25rem",
  },
});

export const cmHighlightStyle = HighlightStyle.define([
  { tag: t.keyword, color: chalky },
  { tag: t.attributeValue, color: cyan },
  { tag: t.meta, color: violet },
  { tag: t.string, color: coral },
  { tag: t.special(t.string), color: coral },
  { tag: t.content, color: "oklch(var(--bc))" },
  { tag: t.definition(t.propertyName), color: "oklch(var(--s))" },
  { tag: t.labelName, color: malibu },
  { tag: t.typeName, color: sage },
  { tag: t.lineComment, color: "oklch(var(--bc) / 0.5)", fontStyle: "italic" },
  { tag: t.comment, color: "oklch(var(--bc) / 0.5)", fontStyle: "italic" },
  { tag: t.separator, color: "oklch(var(--bc))" },
  { tag: t.punctuation, color: whiskey },
  { tag: t.squareBracket, color: "oklch(var(--wa))" },
  { tag: t.brace, color: whiskey },
]);

export const cmHighlightYAML: Extension = [
  yamlLanguage,
  syntaxHighlighting(cmHighlightStyle),
];
