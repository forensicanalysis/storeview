import Vue from 'vue';
import Vuetify from 'vuetify/lib';

import colors from 'vuetify/lib/util/colors';


Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: colors.lightBlue.darken3,
        secondary: colors.pink.darken2,
        accent: colors.pink.darken2,
      },
    },
  },
});
