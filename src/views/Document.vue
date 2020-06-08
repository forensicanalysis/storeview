<!--
Copyright (c) 2020 Siemens AG

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Author(s): Jonas Plum
-->
<template>
  <div>
    <v-tabs
      v-model="tab"
      icons-and-text
      next-icon="mdi-menu-right"
      prev-icon="mdi-menu-left"
      show-arrows>
      <v-tabs-slider/>
      <v-tab :class="detailsTabMenu"
             v-for="view in views"
             :key="view['title']"
             :href="'#tab-'+_.lowerCase(view['title'])">
        {{view['title']}}
        <v-icon>mdi-{{view['icon']}}</v-icon>
      </v-tab>
    </v-tabs>
    <v-divider></v-divider>
    <v-tabs-items v-model="tab">
      <v-tab-item :value="'tab-raw'">
        <v-card flat class="pa-2" style="overflow: auto; max-height: 100%">
          <vue-code-highlight language="json">{{ JSON.stringify(content, null, 2) }}
          </vue-code-highlight>
        </v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-info'">
        <v-card flat class="pa-2" style="overflow: auto; max-height: 100%">
          <JsonToHtml :json="content" rootClassName="wrapper" :inheritClassName="false"/>
        </v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-text'">
        <v-card flat class="pa-2" style="overflow: auto">
          <pre>{{data}}</pre>
        </v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-hex'" style="overflow: auto">
        <!--v-card flat class="pa-2">{{ hexEncode(data) }}</v-card-->
      </v-tab-item>
      <v-tab-item :value="'tab-pdf'" style="overflow: auto">
        <pdf
          ref="pdf"
          :src="data"
          :page="page"
          style="border: 1px solid red"
          @num-pages="numPages = $event"
          @loaded="loaded = true"
        />
      </v-tab-item>
      <v-tab-item :value="'tab-image'" style="overflow: auto">
        <v-card flat class="pa-2" style="overflow: auto; max-height: 100%">
          <img v-if="'export_path' in content"
               :src="'http://localhost:8081/api/file?path='+content.export_path"
               style="width: 100%"/>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import pdf from 'vue-pdf';
import { component as VueCodeHighlight } from 'vue-code-highlight';
import JsonToHtml from '@/components/json-to-html.vue';
import Vuetify from '../plugins/vuetify';

export default {
  name: 'item',
  components: {
    JsonToHtml,
    pdf,
    VueCodeHighlight,
  },
  data() {
    return {
      tab: null,
      data: '',
      active: 'Info',
      views: [
        {
          title: 'Info',
          icon: 'information',
        },
        {
          title: 'Raw',
          icon: 'json',
        },
      ],
      numPages: 1,
      page: 1,
      loaded: false,
    };
  },
  computed: {
    detailsTabMenu() {
      return {
        'detailsTabMenuDark': Vuetify.framework.theme.dark,
        'detailsTabMenuLight': !Vuetify.framework.theme.dark
      }
    },
    theme() {
      return (Vuetify.framework.theme.dark) ? 'dark' : 'light';
    },
    dark() {
      return Vuetify.framework.theme.dark;
    },
    light() {
      return !Vuetify.framework.theme.dark;
    },
  },
  watch: {
    content() {
      if ('export_path' in this.content && this._.endsWith(this.content.export_path, '.pdf')) {
        this.views = this._.filter(this.views, o => o.title !== 'PDF');
        this.views.push({
          title: 'PDF',
          icon: 'file-pdf-outline',
        });
      } else {
        this.views = this._.filter(this.views, o => o.title !== 'PDF');
      }

      if ('export_path' in this.content && this._.endsWith(this.content.export_path, '.jpg')) {
        this.views = this._.filter(this.views, o => o.title !== 'Image');
        this.views.push({
          title: 'Image',
          icon: 'file-image-outline',
        });
      } else {
        this.views = this._.filter(this.views, o => o.title !== 'Image');
      }

      if ('export_path' in this.content) {
        this.views = this._.filter(this.views, o => o.title !== 'Hex' && o.title !== 'Text');
        this.views.push({
          title: 'Hex',
          icon: 'file-document-outline',
        }, {
          title: 'Text',
          icon: 'file-document-outline',
        });
        const that = this;
        this.$http.get(`http://localhost:8081/api/file?path=${this.content.export_path}`)
          .then((response) => {
            that.data = response.data;
          });
      } else {
        this.views = this._.filter(this.views, o => o.title !== 'Hex' && o.title !== 'Text');
      }
    },
  },
  props: {
    content: {
      type: Object,
      required: true,
    },
  },
  methods: {
    hexEncode(s) {
      let hex;
      let i;

      let result = '';
      for (i = 0; i < s.length; i += 1) {
        hex = s.charCodeAt(i)
          .toString(16);
        result += (` ${hex}`).slice(-4);
      }

      return result;
    },
  },
};
</script>

<style lang="sass">

  @import '~animate.css'
  @import '../styles/colors.scss'
  @import '../styles/animations.scss'

  .v-slide-group.v-item-group > .v-slide-group__next, .v-slide-group.v-item-group > .v-slide-group__prev
    min-width: 0
    .theme--light.v-icon
      &:hover
        color: $c-pink

  .detailsTabMenuLight
    transition: $transition-fastest
    &:hover
      font-size: 18px
      color: $c-shadow !important

  .detailsTabMenuDark
    transition: $transition-fastest
    &:hover
      font-size: 18px
      color: $c-pink !important

</style>
