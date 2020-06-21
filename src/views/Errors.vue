<template>
  <v-data-table
    :headers="headers"
    :items="errors"
    :loading="loading"
    :fixed-header="true"
    :footer-props="{'items-per-page-options': [50, 100]}"
    :items-per-page="25"
    style="overflow: visible !important;"
    dense
  >
    <template v-slot:item.title="{ item }">
      <div>{{ title(item) }}</div>
    </template>
    <template v-slot:item.description="{ item }">
      <div>{{ description(item) }}</div>
    </template>
  </v-data-table>
</template>

<script>
  import {invoke} from "../store/invoke";

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

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

        await pause(1500)

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
