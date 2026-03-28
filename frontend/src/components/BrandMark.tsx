export function BrandMark() {
  return (
    <svg className="brand-mark-svg" viewBox="0 0 48 48" aria-hidden="true">
      <defs>
        <linearGradient id="brand-gradient" x1="8" y1="6" x2="38" y2="42" gradientUnits="userSpaceOnUse">
          <stop stopColor="#C8602F" />
          <stop offset="1" stopColor="#27594C" />
        </linearGradient>
      </defs>
      <rect x="4" y="4" width="40" height="40" rx="14" fill="url(#brand-gradient)" />
      <path
        d="M15 15H23C30.7 15 35 19.1 35 24C35 28.9 30.7 33 23 33H15V15ZM22.6 19.6H20.1V28.4H22.6C27 28.4 29.6 26.8 29.6 24C29.6 21.2 27 19.6 22.6 19.6Z"
        fill="#FFF9F1"
      />
      <circle cx="33.5" cy="15.5" r="2.5" fill="#FFF9F1" fillOpacity="0.92" />
    </svg>
  );
}
