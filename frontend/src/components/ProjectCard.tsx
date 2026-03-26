import type { ProjectPayload } from "../types";

export function ProjectCard({ project }: { project: ProjectPayload }) {
  return (
    <article className={`project-card accent-${project.accent}`}>
      <div className="project-heading">
        <div>
          <p className="eyebrow">{project.status}</p>
          <h3>{project.name}</h3>
        </div>
        <a href={project.link} target="_blank" rel="noreferrer">
          访问
        </a>
      </div>
      <p>{project.summary}</p>
      <div className="tag-row">
        {project.techStack.map((item) => (
          <span key={item} className="tag-chip">
            {item}
          </span>
        ))}
      </div>
    </article>
  );
}
