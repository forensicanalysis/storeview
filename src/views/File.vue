<template>
  <div>
    <input type="text" v-model="path" style="width: 100%;">
    <v-tabs v-model="tab" icons-and-text>
      <v-tabs-slider/>
      <v-tab v-if="_.endsWith(path, '.pdf')" href="#tab-pdf">PDF
        <v-icon v-text="'mdi-file-pdf-outline'"/>
      </v-tab>
      <v-tab v-if="isImagePath(path)" href="#tab-image">Image
        <v-icon v-text="'mdi-file-image-outline'"/>
      </v-tab>
      <v-tab v-if="isTextPath(path)" href="#tab-text">Text
        <v-icon v-text="'mdi-file-document-outline'"/>
      </v-tab>
      <v-tab href="#tab-hex">Hex
        <v-icon v-text="'mdi-file-document-outline'"/>
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item :value="'tab-text'">
        <v-card flat class="pa-2" style="overflow: auto">
          <pre>{{data}}</pre>
        </v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-hex'" style="overflow: auto">
        <v-card flat class="pa-2">{{ hex }}</v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-pdf'" style="overflow: auto">
        <pdf
          ref="pdf"
          :src="data"
          :page="1"
          style="border: 1px solid red"
          @num-pages="numPages = $event"
          @loaded="loaded = true"
        />
      </v-tab-item>
      <v-tab-item :value="'tab-image'" style="overflow: auto">
        <v-card flat class="pa-2" style="overflow: auto; max-height: 100%">
          <img
            :src="'http://localhost:8081/api/file?path='+path"
            style="width: 100%" alt=""/>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
  import JsonToHtml from '@/components/json-to-html';
  import pdf from 'vue-pdf';
  import { component as VueCodeHighlight } from 'vue-code-highlight';
  // import hexyjs from 'hexyjs';

  export default {
    name: 'file',
    components: {
      JsonToHtml,
      pdf,
      VueCodeHighlight,
    },
    data: function () {
      return {
        'path': '',
        'tab': null,
        'data': null,
        'hex': '0000',
        'active': 'Info',
        'loaded': false,
      };
    },
    watch: {
      path: function () {
        let vm = this;
        this.$http.get('http://localhost:8081/api/file?path=' + this.path)
          .then(response => {
            // vm.hex = hexyjs.strToHex(response.data);
            vm.data = response.data;
          });
      }
    },
    props: {
      /* future_path: {
        type: Object,
        required: true,
      }, */
    },
    mounted() {
      this.path = 'WindowsRecycleBin/$I9F219A.jpg';
    },
    methods: {
      hasEnding(path, endings) {
        for (let i = 0; i < endings.length; i++) {
          if (_.endsWith(path, '.' + endings[i])) {
            return true;
          }
        }
        return false;
      },
      isImagePath(path) {
        return this.hasEnding(path, ['jpg', 'jpeg', 'png', 'svg']);
      },
      isTextPath(path) {
        return this.hasEnding(path, ['txt', 'md']);
      }
    }
  };
</script>

