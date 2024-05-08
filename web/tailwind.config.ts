import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
    },
  },
  daisyui: {
    themes: [
      {
        mytheme: {
          "primary": "#00c7f5",
          "primary-content": "#000e15",
          "secondary": "#00bb63",
          "secondary-content": "#000d03",
          "accent": "#e67600",
          "accent-content": "#130500",
          "neutral": "#362913",
          "neutral-content": "#d3d0cb",
          "base-100": "#2f2032",
          "base-200": "#271a2a",
          "base-300": "#201522",
          "base-content": "#d1cdd2",
          "info": "#00a6c3",
          "info-content": "#000a0e",
          "success": "#86b500",
          "success-content": "#060c00",
          "warning": "#ffa63e",
          "warning-content": "#160a01",
          "error": "#ff2c4a",
          "error-content": "#160102",
        },
      },
    ],
  },
  plugins: [require("daisyui")],
};
export default config;
