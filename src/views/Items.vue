<template>
  <v-row class="mt-2" style="flex-wrap: nowrap;">
    <v-col
      ref="drawerLeft"
      :style="'overflow: hidden; width: '+navigationLeft.width+'px; flex-basis: '+navigationLeft.width+'px'"
      class="flex-grow-0 flex-shrink-0 ml-3"
    >
      <v-card>
        <v-list dense>
          <v-subheader>Search</v-subheader>
          <v-text-field
            v-model="search"
            prepend-icon="mdi-magnify"
            clear-icon="mdi-close-circle"
            clearable
            @keyup.enter.native="searchFilter"
            @click:append="searchFilter"
            @click:clear="clearFilter"
            class="mx-3 my-0 py-0"
          ></v-text-field>
        </v-list>
      </v-card>
      <v-card class="mt-2">
        <v-list dense>
          <v-subheader>Type</v-subheader>
          <v-list-item-group v-model="itemTypeIndex" color="primary">
            <v-list-item
              v-for="(table, i) in _.sortBy($store.state.tables, ['name'])"
              :key="i"
              @click="loadList(table['name'])"
              class="px-4"
            >
              <v-list-item-icon>
                <v-icon v-text="'mdi-'+$store.state.templates[table['name']].icon"/>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title v-text="_.startCase(table.name)"/>
              </v-list-item-content>
              <v-list-item-icon>
                <span v-if="'count' in table">{{table['count']}}</span>
              </v-list-item-icon>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-card>
      <v-card class="mt-2">
        <v-list dense>
          <v-list-item-group>
            <v-subheader>Location</v-subheader>
            <v-treeview
              v-if="refresh && directories.length > 0"
              activatable
              hoverable
              dense
              transition
              @update:active="updatedir"
              :active.sync="active"
              :items="directories"
              :load-children="fetchTreeChildren"
              :open.sync="open"
              item-key="path"
              color="primary"
              ref="treeView"
            >
              <template v-slot:prepend="{ item, open }">
                <v-icon v-if="item.children">
                  {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                </v-icon>
                <v-icon v-if="!item.children">mdi-folder-outline</v-icon>
              </template>
            </v-treeview>
          </v-list-item-group>
        </v-list>
      </v-card>
      <v-card class="mt-2">
        <v-spacer></v-spacer>
        <v-btn
          small
          icon
          right
          @click="toogleLeft"
        >
          <v-icon>mdi-menu</v-icon>
        </v-btn>
      </v-card>
    </v-col>
    <v-col class="flex-grow-1 flex-shrink-1">
      <v-card class="pt-2">
        <v-data-table
          :headers="$store.state.headers"
          :items="$store.state.items"
          :loading="loading"
          :options.sync="options"
          :server-items-length="$store.state.itemCount"
          @update:options="updateopt"
          :fixed-header="true"
          multi-sort
          @click:row="select"
          :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
          show-select
          dense
        >
          <template v-slot:body.prepend>
            <tr>
              <td></td>
              <td v-for="h in $store.state.headers" :key="h.text" role="columnheader" scope="col">
                <v-text-field
                  v-model="itemscol[h.value]"
                  @keyup.enter.native="searchFilter"
                  hide-details
                  label="Filter"
                  dense
                />
              </td>
            </tr>
          </template>

          <template v-slot:item.atime="{ item }">
            <div>{{ toLocal(item.atime) }}</div>
          </template>
          <template v-slot:item.mtime="{ item }">
            <div>{{ toLocal(item.mtime) }}</div>
          </template>
          <template v-slot:item.ctime="{ item }">
            <div>{{ toLocal(item.ctime) }}</div>
          </template>
          <template v-slot:item.origin="{ item }">
            <div><span v-if="'path' in item.origin">{{item.origin.path}}</span></div>
          </template>
          <template v-slot:item.size="{ item }">
            <div>{{ humanBytes(item.size, true) }}</div>
          </template>
        </v-data-table>
      </v-card>
    </v-col>
    <v-col
      :style="'overflow: hidden; width: '+navigationRight.width+'px; flex-basis: '+navigationRight.width+'px'"
      ref="drawerRight"
      class="flex-grow-0 flex-shrink-0 mr-3"
    >
      <v-card>
        <v-toolbar
          v-if="!_.isEmpty($store.state.item) && !navigationRight.mini"
          class="elevation-0"
          dense>
          <v-spacer></v-spacer>
          <v-btn
            small
            icon
            @click="toogleRight"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-toolbar>
        <transition name="fade">
          <item v-show="!navigationRight.mini" :content="$store.state.item"/>
        </transition>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
  import {DateTime} from 'luxon';
  import item from '@/views/Document.vue';

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

  export default {

    name: 'items',

    components: {
      item,
    },

    data() {
      return {

        itemscol: {},

        filter: {
          table: this.$store.state.type,
          columns: {},
        },

        search: '',
        active: [],
        open: [],
        refresh: true,

        navigationLeft: {
          mini: false,
          shown: true,
          width: 250,
          oldWidth: 250,
          borderSize: 3,
        },

        navigationRight: {
          mini: true,
          shown: true,
          width: 56,
          oldWidth: 500,
          borderSize: 3,
        },

        itemTypeIndex: 0,
        loading: false,
        options: {},

        directories: []
      };
    },

    computed: {

      count() {
        for (let i = 0; i < this.$store.state.tables.length; i += 1) {
          if (this.$store.state.type === this.$store.state.tables[i].name) {
            console.log(this.$store.state.tables[i]);
            return this.$store.state.tables[i].count;
          }
        }
        return 0;
      },

      paneSize: {
        get() {
          return this.$store.state.listPane;
        },
        set(value) {
          this.$store.commit('setlistPane', value);
        },
      },

    },

    methods: {

      toggle(navigation) {
        if (!navigation.mini) {
          navigation.oldWidth = parseInt(navigation.width, 10);
          navigation.width = 80;
        }
        navigation.mini = !navigation.mini;

        if (!navigation.mini) {
          navigation.width = parseInt(navigation.oldWidth, 10);
        }
      },

      toogleLeft() {
        this.toggle(this.navigationLeft);
      },

      toogleRight() {
        this.toggle(this.navigationRight);
      },

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

      forceRefresh() {
        this.refresh = false;
        this.$nextTick(() => {
          this.refresh = true;
        });
      },

      async getDirectories() {
        var that = this;
        if (this.directories.length === 0) {

          this.$store.dispatch('loadDirectories', {path: ''})
            .then((response) => {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  that.directories.push(response[i]);
                }
              }
            })
            .catch(error => console.warn(error));

          console.log('INITIAL: ', JSON.stringify(this.directories, null, 1));
        }
      },

      async loadList(table) {
        this.directories = [];

        this.itemscol = {};

        this.active = [];
        this.open = [];

        this.$store.commit('setItem', {});
        if (!this.navigationRight.mini) {
          this.toogleRight();
        }
        this.$store.commit('setTable', table);
        this.$store.dispatch('loadItems')
          .then(() => {
            for (let i = 0; i < this.$store.state.headers.length; i += 1) {
              this.itemscol[this.$store.state.headers[i].value] = "";
            }

            this.getDirectories();
            this.forceRefresh();

            this.filter = {
              table: this.$store.state.type,
              columns: {},
            };
          });
      },

      select(e) {
        this.$store.commit('setItem', e);
        if (this.navigationRight.mini) {
          this.toogleRight();
        }
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      updateopt(e) {
        console.log('updateopt', e);
        this.$store.commit('setOffset', (e.page - 1) * e.itemsPerPage);
        this.$store.commit('setLimit', e.itemsPerPage);
        const sort = {
          table: this.$store.state.type,
          columns: {},
        };
        for (let i = 0; i < e.sortBy.length; i += 1) {
          if (e.sortDesc[i]) {
            sort.columns[e.sortBy[i]] = 'DESC';
          } else {
            sort.columns[e.sortBy[i]] = 'ASC';
          }
        }

        this.$store.commit('setSort', sort);
        this.$store.dispatch('loadItems');
      },

      updatedir(e) {
        console.log('updatedir', e);
        console.log(JSON.stringify(this.directories, null, 1));

        let slash = '';
        let column = '';
        const {type} = this.$store.state;

        if (type === 'directory') {
          console.log('DIR');
          slash = '/';
          column = 'path';
        } else if (type === 'file') {
          console.log('FILE');
          slash = '/';
          column = 'origin.path';
        } else if (type === 'windows-registry-key') {
          console.log('KEY');
          slash = '\\';
          column = 'key';
        } else {
          console.log('TABLE DOES NOT EXIST!');
        }

        if (e.length === 0) {
          this.filter.table = type;
          this.filter.columns[column] = [];
        } else {
          this.filter.table = type;
          this.filter.columns[column] = [e[0] + slash];
        }

        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');
      },

      searchFilter() {
        this.filter.table = this.$store.state.type;

        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          if (key === 'origin') {
            column = 'origin.path';
          } else {
            column = key;
          }
          this.filter.columns[column] = [value];
        }

        this.filter.columns['elements'] = [this.search];
        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');
      },

      clearFilter() {
        this.search = '';
        this.searchFilter();
      },

      async fetchTreeChildren(item) {
        const directories = [];

        this.$store.dispatch('loadDirectories', {path: item.path})
          .then((response) => {
            if ((response.length === 1) && (response[0].name === '/')) {
              delete item.children;
            } else {
              for (let i = 0; i < response.length; i += 1) {
                if (response[i].name !== '/') {
                  directories.push(response[i]);
                }
              }
            }
          })
          .catch(error => console.warn(error));

        await pause(1000);

        const key = item.path;
        const parentNode = this.$refs.treeView.nodes[key];

        let childNode;
        directories.forEach((child) => {
          childNode = {
            ...parentNode,
            item: child,
            vnode: null,
          };
          this.$refs.treeView.nodes[child.path] = childNode;
        });

        for (let i = 0; i < directories.length; i += 1) {
          item.children.push(directories[i]);
        }
      },

    },

    mounted() {
      this.getDirectories();
    },

    destroyed() {
      this.$store.commit('setItem', {});
    },

  };


</script>

<style>
  table tr td {
    cursor: pointer !important;
  }

  * {
    border-radius: 0 !important;
  }

  .fade-enter-active, .fade-leave-active {
    transition: opacity .2s;
  }

  .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */
  {
    opacity: 0;
  }

  .v-treeview-node__label, .v-data-table .v-text-field, .v-data-table .v-label {
    font-size: 0.8125rem;
  }

  .v-data-table .v-label--active {
    opacity: 0;
  }

  .v-data-table .v-text-field > .v-input__control > .v-input__slot:before {
    border-style: none;
  }
</style>
