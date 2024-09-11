/** @type {import('tailwindcss').Config} */
export const content = ["./view/**/*.templ"];
export const theme = {
  extend: {
    fontFamily: {
      sans: ["Roboto", "sans-serif"],
      mono: ["Ubuntu Mono", "monospace"],
    },
  },
};
export const plugins = [require("@tailwindcss/forms"), require("daisyui")];
export const daisyui = {
  themes: [
    {
      light: {
        ...require("daisyui/src/theming/themes")["light"],
        primary: "#5f3dc4",
        "primary-content": "white",
        secondary: "#1864ab",
        "secondary-content": "white",
        accent: "#f6f8fa",
        "success-content": "white",
        "error-content": "white",
        "--rounded-btn": "0.5rem",
        "--rounded-box": "0.5rem",
      },
    },
  ],
};
