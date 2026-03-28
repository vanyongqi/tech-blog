import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { BrandMark } from "./BrandMark";

describe("BrandMark", () => {
  it("renders branded svg mark", () => {
    const html = renderToStaticMarkup(<BrandMark />);

    expect(html).toContain("brand-mark-svg");
    expect(html).toContain("<svg");
    expect(html).toContain("brand-gradient");
  });
});
