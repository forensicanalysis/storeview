<template>
  <v-container fluid>
    <v-card>
      <v-breadcrumbs :items="breadcrumbs">
        <v-breadcrumbs-item
          slot="item"
          slot-scope="{ item }"
          exact
          @click="loadFiles( item.path )">
          {{ item.text }}
        </v-breadcrumbs-item>
      </v-breadcrumbs>
      <v-data-table
        :headers="headers"
        :items="files"
        :fixed-header="true"
        :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
        style="overflow: visible !important;"
        dense
      >
        <template v-slot:item.name="{ item }">
          <a v-if="item.dir" @click="loadFiles(item.path)">{{ item.name }}</a>
          <div v-else>{{ item.name }}</div>
        </template>
        <template v-slot:item.mtime="{ item }">
          <div>{{ toLocal(item.mtime) }}</div>
        </template>
        <template v-slot:item.size="{ item }">
          <div>{{ humanBytes(item.size, true) }}</div>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-icon
            v-if="!item.dir"
            small
            @click="editItem(item)"
          >
            mdi-download
          </v-icon>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>

<script>
  import {invoke} from "../store/invoke";
  import {DateTime} from "luxon";

  export default {
    name: 'files',
    data() {
      return {
        breadcrumbs: [{text: '/', path: '/', link: true}],
        headers: [
          {text: "Name", value: "name"},
          // {text: "Path", value: "path"},
          {text: "Modified", value: "mtime"},
          {text: "Size", value: "size"},
          {text: "Download", value: "actions", sortable: false},
        ],
        files: [],
      };
    },

    methods: {
      humanBytes(bytes, si) {
        const thresh = si ? 1000 : 1024;
        if (Math.abs(bytes) < thresh) {
          return `${bytes}B`;
        }
        const units = si
          ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        do {
          bytes /= thresh;
          u += 1;
        } while (Math.abs(bytes) >= thresh && u < units.length - 1);
        return `${bytes.toFixed(1)}${units[u]}`;
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      loadFiles(path) {
        let breadcrumbs = [];
        let parts = path.split("/");
        for (let i = 0; i < parts.length; i += 1) {
          breadcrumbs.push({
            text: parts[i],
            href: parts[i],
            path: parts.slice(0, i + 1).join("/"),
            link: true,
          })
        }
        this.breadcrumbs = breadcrumbs;
        invoke('GET', '/files?path=' + path, [], (data) => {
          console.log(data);
          this.files = data.elements;
        });
      },
    },

    mounted() {
      this.loadFiles("/");
    },
  };
</script>
