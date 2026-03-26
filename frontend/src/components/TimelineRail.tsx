import type { TimelineEntryPayload } from "../types";

export function TimelineRail({
  entries,
}: {
  entries: TimelineEntryPayload[];
}) {
  return (
    <div className="timeline-rail">
      {entries.map((entry) => (
        <article key={`${entry.period}-${entry.title}`} className="timeline-card">
          <div className="timeline-period">{entry.period}</div>
          <div>
            <h3>{entry.title}</h3>
            <p>{entry.description}</p>
          </div>
        </article>
      ))}
    </div>
  );
}
