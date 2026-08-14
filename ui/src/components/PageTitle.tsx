export interface PageTitleProps {
  title: string;
  description: string;
}

export function PageTitle({ title, description }: PageTitleProps) {
  return (
    <div className="mb-10">
      <h1 className="text-3xl font-semibold">{title}</h1>
      <p className="text-base-content/60">{description}</p>
    </div>
  );
}
