export type MarkdownImageSize = "sm" | "md" | "lg" | "full";
export type MarkdownImageAlign = "left" | "center" | "right";

export type MarkdownImageOptions = {
  size: MarkdownImageSize;
  align: MarkdownImageAlign;
  width?: string;
};

const defaultOptions: MarkdownImageOptions = {
  size: "lg",
  align: "center",
};

export function parseMarkdownImageOptions(title?: string | null): MarkdownImageOptions {
  if (!title) {
    return defaultOptions;
  }

  const options: MarkdownImageOptions = { ...defaultOptions };
  for (const token of title.split(/\s+/).filter(Boolean)) {
    const [rawKey, rawValue] = token.split("=");
    const key = rawKey?.trim().toLowerCase();
    const value = rawValue?.trim().toLowerCase();

    if (!key || !value) {
      continue;
    }

    if (key === "size" && isImageSize(value)) {
      options.size = value;
      continue;
    }

    if (key === "align" && isImageAlign(value)) {
      options.align = value;
      continue;
    }

    if (key === "width" && isWidthValue(value)) {
      options.width = normalizeWidthValue(value);
    }
  }

  return options;
}

export function buildMarkdownImageSnippet(alt: string, url: string, options: Partial<MarkdownImageOptions> = {}) {
  const merged = { ...defaultOptions, ...options };
  const tokens = [`size=${merged.size}`, `align=${merged.align}`];
  if (merged.width) {
    tokens.push(`width=${normalizeWidthValue(merged.width)}`);
  }
  return `![${alt}](${url} "${tokens.join(" ")}")`;
}

function isImageSize(value: string): value is MarkdownImageSize {
  return value === "sm" || value === "md" || value === "lg" || value === "full";
}

function isImageAlign(value: string): value is MarkdownImageAlign {
  return value === "left" || value === "center" || value === "right";
}

function isWidthValue(value: string) {
  return /^\d+$/.test(value) || /^\d+%$/.test(value);
}

function normalizeWidthValue(value: string) {
  return /^\d+$/.test(value) ? `${value}px` : value;
}
