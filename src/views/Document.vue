<template>
  <div>
    <v-tabs v-model="tab" icons-and-text>
      <v-tabs-slider/>
      <v-tab v-for="view in views" :key="view['title']" :href="'#tab-'+_.lowerCase(view['title'])">
        {{view['title']}}
        <v-icon>mdi-{{view['icon']}}</v-icon>
      </v-tab>
    </v-tabs>
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
  import JsonToHtml from '@/components/json-to-html';
  import pdf from 'vue-pdf';
  import { component as VueCodeHighlight } from 'vue-code-highlight';

  export default {
    name: 'item',
    components: {
      JsonToHtml,
      pdf,
      VueCodeHighlight
    },
    data: function () {
      return {
        'tab': null,
        'data': '',
        'active': 'Info',
        'views': [
          {
            'title': 'Info',
            'icon': 'information'
          },
          {
            'title': 'Raw',
            'icon': 'json'
          },
        ],
        'numPages': 1,
        'page': 1,
        'loaded': false
      };
    },
    watch: {
      content: function () {
        if ('export_path' in this.content && this._.endsWith(this.content['export_path'], '.pdf')) {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'PDF';
          });
          this.views.push({
            'title': 'PDF',
            'icon': 'file-pdf-outline'
          });
        } else {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'PDF';
          });
        }

        if ('export_path' in this.content && this._.endsWith(this.content['export_path'], '.jpg')) {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'Image';
          });
          this.views.push({
            'title': 'Image',
            'icon': 'file-image-outline'
          });
        } else {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'Image';
          });
        }

        if ('export_path' in this.content) {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'Hex' && o['title'] !== 'Text';
          });
          this.views.push({
            'title': 'Hex',
            'icon': 'file-document-outline'
          }, {
            'title': 'Text',
            'icon': 'file-document-outline'
          });
          let that = this;
          this.$http.get('http://localhost:8081/api/file?path=' + this.content['export_path'])
            .then(response => {
              that.data = response.data;
            });
        } else {
          this.views = this._.filter(this.views, function (o) {
            return o['title'] !== 'Hex' && o['title'] !== 'Text';
          });
        }
      }
    },
    props: {
      content: {
        type: Object,
        required: true,
      },
    },
    methods: {
      hexEncode: function (s) {
        let hex,
          i;

        let result = '';
        for (i = 0; i < s.length; i++) {
          hex = s.charCodeAt(i)
            .toString(16);
          result += (' ' + hex).slice(-4);
        }

        return result;
      }
    }
  };
</script>

