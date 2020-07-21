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
  <div class="scrollableArea">
    <v-data-table
      :headers="headers"
      :items="errors"
      :loading="loading"
      :fixed-header="true"
      :footer-props="{'items-per-page-options': [10, 25, 50, 1000]}"
      :items-per-page="25"
      dense
    >
      <template v-slot:item.title="{ item }">
        <div>{{ title(item) }}</div>
      </template>
      <template v-slot:item.description="{ item }">
        <div>{{ description(item) }}</div>
      </template>
    </v-data-table>
  </div>
</template>

<script>
  import {invoke} from "../store/invoke";

  export default {
    name: 'errors',
    data() {
      return {
        loading: true,
        headers: [
          // {text: "ID", value: "id"},
          {text: "Title", value: "title"},
          {text: "Type", value: "type"},
          // {text: "Description", value: "description"},
          {text: "Errors", value: "errors"},
        ],
        errors: [],
      };
    },

    methods: {
      async loadFiles() {

        invoke('GET', '/errors', [], (data) => {
          this.errors = data.elements;
        }).then(() => {
          this.loading = false;
        });
      },

      title(element) {
        if (_.has(this.$store.state.templates, element['type'])) {
          return element[this.$store.state.templates[element['type']].headers[0].value]
        }
        if (_.has(element, 'name')) {
          return element['name'];
        }
        if (_.has(element, 'key')) {
          return element['key'];
        }
        if (_.has(element, 'title')) {
          return element['title'];
        }
        return "";
      },
      description(element) {
        if (_.has(element, 'description')) {
          return element['name'];
        }
        if (_.has(element, 'desc')) {
          return element['key'];
        }
        return "";
      },
    },

    mounted() {
      this.loadFiles();
    },
  };
</script>
