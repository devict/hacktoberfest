const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin')

module.exports = {
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  experimental: {
    uniformColorPalette: true,
    defaultLineHeights: true,
    extendedFontSizeScale: true,
  },
  purge: [
    './templates/**/*.tmpl',
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
      fontSize: {
        '10xl': '10rem',
      },
      height: {
        '3': '.7rem',
        '1/3': '33.333%',
        '1/2': '50%',
      },
      maxHeight: {
        '10': '2.5rem',
      },
    },
  },
  variants: {
    textColor: ({ before }) => before(['group-hover'], 'hover'),
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/custom-forms'),
    plugin(function({ addComponents, theme }) {
      addComponents({
        '.paging_simple_numbers': {
          borderRadius: theme('borderRadius.default'),
          overflow: 'hidden',
          border: `1px solid ${theme('colors.gray.200')}`,
          borderRight: 'none',
          display: 'inline-flex',
          '& > span': {
            display: 'flex',
          },
          '& a': {
            display: 'block',
            padding: '.5rem .75rem',
            margin: '-1px 0 -1px -1px',
            cursor: 'pointer',
            border: `1px solid ${theme('colors.gray.200')}`,
          },
        },
      })
    })
  ],
}
