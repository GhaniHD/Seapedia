// NOTE: component/file names kept as-is (WaveDivider, ShipWheelMark) even
// though the visuals are no longer nautical, so every existing import across
// the app keeps working unchanged after this redesign to a marketplace look.

export function WaveDivider() {
  // A plain thin divider — the nautical wave shape is gone in the
  // marketplace redesign, but the component is kept as a no-op so call
  // sites don't need to change.
  return null;
}

export function ShipWheelMark({ className = "h-8 w-8" }: { className?: string }) {
  // Simple shopping-bag glyph used as the SEAPEDIA logo mark.
  return (
    <svg viewBox="0 0 24 24" className={className} fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <path
        d="M6 8h12l-1 12a2 2 0 0 1-2 2H9a2 2 0 0 1-2-2L6 8Z"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinejoin="round"
      />
      <path d="M9 8V6a3 3 0 1 1 6 0v2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}
