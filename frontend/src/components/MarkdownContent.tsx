import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { parseMarkdownImageOptions } from "../lib/markdownImages";

export function MarkdownContent({ content }: { content: string }) {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={{
        a: ({ ...props }) => <a {...props} className="inline-link" target="_blank" rel="noreferrer" />,
        img: ({ ...props }) => {
          const options = parseMarkdownImageOptions(props.title);
          return (
            <img
              {...props}
              alt={props.alt ?? ""}
              loading="lazy"
              className={`markdown-image markdown-image-size-${options.size} markdown-image-align-${options.align}`}
              style={options.width ? { width: options.width, maxWidth: "100%" } : undefined}
            />
          );
        },
      }}
    >
      {content}
    </ReactMarkdown>
  );
}
