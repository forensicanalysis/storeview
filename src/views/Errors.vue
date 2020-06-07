<template>
  <v-data-table
    :headers="headers"
    :items="errors"
    :fixed-header="true"
    :footer-props="{'items-per-page-options': [50, 100]}"
    :items-per-page="25"
    style="overflow: visible !important;"
    dense
  >
  </v-data-table>
</template>

<script>
  import {invoke} from "../store/invoke";

  export default {
    name: 'logs',
    data() {
      return {
        headers: [
          {text: "Errors", value: "errors"},
          {text: "ID", value: "id"},
        ],
        errors: [],
      };
    },

    methods: {
      loadFiles() {
        invoke('GET', '/errors', [], (data) => {
          this.errors = data.elements;
        });
      },
    },

    mounted() {
      this.loadFiles();
    },
  };
</script>
