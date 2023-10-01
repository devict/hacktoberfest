const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin')

module.exports = {
  content: [
    './templates/**/*.tmpl',
    './public/js/**/*.js',
  ],
  theme: {
    typography: (theme) => ({
      default: {
        css: {
          a: {
            color: theme('colors.teal.600'),
            '&:hover': {
              color: theme('colors.teal.700'),
            },
          },
          'ul > li::before': {
            backgroundColor: theme('colors.teal.600'),
          },
        },
      },
    }),
    extend: {
      colors: {
        teal: {
          50: '#F9FEFF',
          100: '#C0EAF4',
          200: '#98D7E6',
          300: '#7DC3D3',
          400: '#6FB1C4',
          500: '#6199AF',
          600: '#507B92',
          700: '#3C576C',
          800: '#233040',
          900: '#0A0A12',
        },
        orange: {
          ...defaultTheme.colors.orange,
          400: '#F3A856',
          500: '#F7941E',
          600: '#ae641b',
          700: '#A55302',
          800: '#944a01',
        }
      },
      fontFamily: {
        'sans': ['Montserrat', ...defaultTheme.fontFamily.sans],
        'display': ['Germania\\ One', 'cursive'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
}
