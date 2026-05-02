export const ComingSoon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 800 400"
      width="100%"
      height="100%"
      style={{ display: "block" }}
    >
      {/* Icon: clock/construction */}
      <rect
        x="375"
        y="100"
        width="50"
        height="50"
        rx="4"
        fill="none"
        stroke="#94a3b8"
        strokeWidth="2"
      />
      <line
        x1="400"
        y1="115"
        x2="400"
        y2="130"
        stroke="#94a3b8"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <line
        x1="400"
        y1="130"
        x2="410"
        y2="140"
        stroke="#94a3b8"
        strokeWidth="2"
        strokeLinecap="round"
      />

      {/* Heading */}
      <text
        x="400"
        y="200"
        textAnchor="middle"
        fontFamily="system-ui, sans-serif"
        fontSize="22"
        fontWeight="600"
        fill="#1e293b"
      >
        Coming Soon
      </text>

      {/* Subtext */}
      <text
        x="400"
        y="230"
        textAnchor="middle"
        fontFamily="system-ui, sans-serif"
        fontSize="14"
        fill="#64748b"
      >
        This section is under development.
      </text>
      <line
        x1="340"
        y1="255"
        x2="460"
        y2="255"
        stroke="#e2e8f0"
        strokeWidth="1"
      />
      <rect x="300" y="268" width="200" height="22" rx="11" fill="#f1f5f9" />
      <text
        x="400"
        y="283"
        textAnchor="middle"
        fontFamily="system-ui, sans-serif"
        fontSize="11"
        fill="#94a3b8"
        letterSpacing="1"
      >
        IN PROGRESS
      </text>
    </svg>
  );
};
