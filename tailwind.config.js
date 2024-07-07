/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./cmd/web/**/*.html', './cmd/web/**/*.templ'],
  theme: {
    extend: {
      colors: {
        blue: { primary: '#209cee', hover: '#108de0' },
        green: { primary: '#22c55e' },
      },
      fontFamily: {
        primary: ['"JetBrains Mono"'],
        secondary: ['"Press Start 2P"'],
      },
      fontSize: {
        xxs: ['0.6rem', { lineHeight: '0.9rem' }],
      },
    },
  },
  plugins: [],
};
