import { createTheme } from "@mui/material";
const palette = {
  primary: {
    main: "hsl(260, 71%, 50%)",
    light: "hsl(260, 71%, 60%)",
    dark: "hsl(260, 71%, 40%)",
  },
  secondary: {
    main: "hsl(185, 71%, 50%)",
    light: "hsl(185, 71%, 60%)",
    dark: "hsl(185, 71%, 40%)",
  },
  error: {
    main: "hsl(354, 84%, 51%)",
  },
  warning: {
    main: "hsl(43, 100%, 60%)",
    light: "hsl(43, 100%, 71%)",
    dark: "hsl(29, 100%, 50%)",
  },
  success: {
    main: "hsl(127, 50%, 50%)",
  },
  nudge: {
    main: "hsl(48, 86%, 50%)",
  },
  selected: {
    main: "hsl(260, 23%, 96%)",
    dark: "hsl(260, 30%, 78%)",
  },
};

export const theme = createTheme({
  breakpoints: {
    values: {
      xs: 0,
      sm: 600,
      md: 768,
      lg: 1280,
      xl: 1920,
    },
  },
  palette,
  typography: {
    fontFamily: "sofia-pro,sans-serif",
  },
});
