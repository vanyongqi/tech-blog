import { describe, expect, it } from "vitest";
import { buildMarkdownImageSnippet, parseMarkdownImageOptions } from "./markdownImages";

describe("markdownImages", () => {
  it("parses standard image layout options", () => {
    expect(parseMarkdownImageOptions("size=sm align=right width=320")).toEqual({
      size: "sm",
      align: "right",
      width: "320px",
    });
  });

  it("falls back to unified defaults", () => {
    expect(parseMarkdownImageOptions("")).toEqual({
      size: "lg",
      align: "center",
    });
  });

  it("builds a unified markdown image snippet", () => {
    expect(buildMarkdownImageSnippet("示意图", "/api/assets/1", { size: "full", align: "center" })).toBe(
      '![示意图](/api/assets/1 "size=full align=center")',
    );
  });
});
