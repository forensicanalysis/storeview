<template>
  <div>
    <v-navigation-drawer
      app
      clipped
      stateless
      :width="navigationLeft.width"
      v-model="navigationLeft.shown"
      ref="drawerLeft"
      style="overflow: visible"
    >
      <v-btn
        color="pink"
        fab
        dark
        small
        absolute
        top
        right
        @click="toogleLeft"
        style="margin-top: 25px; margin-right: -35px;"
      >
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-list class="mt-6" dense>
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
          <hr/>
          <h2>Navigation</h2>
          <v-icon
            clickable
            @click="getDirectories"
            color="purple darken-2">
            mdi-menu
          </v-icon>
          <v-treeview
            v-if="refresh && showDirectories"
            activatable
            hoverable
            dense
            transition
            @update:active="updatedir"
            :active.sync="active"
            :items="directories"
            :load-children="fetch"
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
    </v-navigation-drawer>

    <v-container fluid>
      <v-row align="center">

        <v-col class="d-flex" cols="12" sm="9">
          <v-text-field
            style="margin-top: 50px"
            v-model="search"
            :append-icon="'mdi-send'"
            outlined
            clear-icon="mdi-close-circle"
            clearable
            label="Search"
            type="text"
            @click:append="searchFilter"
            @click:clear="clearFilter"
          ></v-text-field>
        </v-col>

        <v-col class="d-flex" cols="12" sm="3">
          <v-select
            v-if="refresh"
            style="margin-top: 50px; width: 300px"
            :items="$store.state.headers"
            v-model="searchField"
            outlined
            label="Field"
          ></v-select>
        </v-col>

      </v-row>
    </v-container>

    <v-card class="mt-10" outlined style="margin-top: 3px !important;">
      <v-data-table
        :headers="$store.state.headers"
        :items="$store.state.items"
        :loading="loading"
        :options.sync="options"
        :server-items-length="count"
        @update:options="updateopt"
        dense
        :fixed-header="true"
        multi-sort
        @click:row="select"
        :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
      >
        <template v-slot:body.prepend>
          <tr>
            <th v-for="h in $store.state.headers" role="columnheader" scope="col">
              <v-combobox
                v-model.sync="itemscol[h.value]"
                :search-input.sync="searchcol[h.value]"
                multiple
                small-chips
                dense
                hide-no-data
                deletable-chips
                clearable
                @change="searchFilter"
                append-icon=""
              />
            </th>
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
    <v-navigation-drawer
      app
      clipped
      right
      stateless
      :width="navigationRight.width"
      v-model="navigationRight.shown"
      ref="drawerRight"
      style="overflow: visible;"
    >
      <v-btn
        v-if="!_.isEmpty($store.state.item)"
        color="pink"
        fab
        dark
        small
        absolute
        top
        left
        @click="toogleRight"
        style="margin-top: 25px; margin-left: -35px;"
      >
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <transition name="fade">
        <item v-show="!navigationRight.mini" :content="$store.state.item"/>
      </transition>
    </v-navigation-drawer>
  </div>
</template>

<script>
  import JsonToHtml from '../components/json-to-html';
  import item from '@/views/Document';
  import { DateTime } from 'luxon';

  const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

  export default {

    name: 'items',

    components: {
      JsonToHtml,
      item
    },

    data: function () {
      return {

        itemscol: {},
        searchcol: {},

        showDirectories: false,

        filter: {
          'table': this.$store.state.type,
          'columns': {}
        },

        searchField: '',
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

      };

    },

    computed: {

      count: function () {
        for (let i = 0; i < this.$store.state.tables.length; i++) {
          if (this.$store.state.type === this.$store.state.tables[i]['name']) {
            console.log(this.$store.state.tables[i]);
            return this.$store.state.tables[i]['count'];
          }
        }
        return 0;
      },

      directionLeft() {
        return this.navigationLeft.shown === false ? 'Open' : 'Closed';
      },

      directionRight() {
        return this.navigationLeft.shown === false ? 'Open' : 'Closed';
      },

      paneSize: {
        get() {
          return this.$store.state.listPane;
        },
        set(value) {
          this.$store.commit('setlistPane', value);
        }
      },

      directories() {
        return [];
      },

    },

    methods: {

      toggle(navigation) {
        if (!navigation.mini) {
          navigation.oldWidth = parseInt(navigation.width, 10);
          navigation.width = 56;
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

      setBorderWidthLeft() {
        let i = this.$refs.drawerLeft.$el.querySelector(
          '.v-navigation-drawer__border'
        );
        i.style.width = this.navigationLeft.borderSize + 'px';
        i.style.cursor = 'ew-resize';
      },

      setBorderWidthRight() {
        let i = this.$refs.drawerRight.$el.querySelector(
          '.v-navigation-drawer__border'
        );
        i.style.width = this.navigationRight.borderSize + 'px';
        i.style.cursor = 'ew-resize';
      },

      setEventsAll() {
        this.setEvents(this.$refs.drawerLeft.$el, this.navigationLeft);
        this.setEvents(this.$refs.drawerRight.$el, this.navigationRight);
      },

      setEvents(el, navigation) {
        const minSize = navigation.borderSize;
        const drawerBorder = el.querySelector('.v-navigation-drawer__border');
        const direction = el.classList.contains('v-navigation-drawer--right')
          ? 'right'
          : 'left';

        function resize(e) {
          document.body.style.cursor = 'ew-resize';
          let f = direction === 'right'
            ? document.body.scrollWidth - e.clientX
            : e.clientX;
          el.style.width = f + 'px';
        }

        drawerBorder.addEventListener(
          'mousedown',
          function (e) {
            if (e.offsetX < minSize) {
              // m_pos = e.x;
              if (navigation.mini) {
                return false;
              }

              el.style.transition = 'initial';
              document.addEventListener('mousemove', resize, false);
            }
          },
          false
        );

        document.addEventListener(
          'mouseup',
          function () {
            el.style.transition = '';
            navigation.width = el.style.width;

            /*  console.log(el.style.width);
            if (parseInt(el.style.width, 10) > 80) {
              console.log('false');
              navigation.mini = false;
            } */

            document.body.style.cursor = '';
            document.removeEventListener('mousemove', resize, false);
          },
          false
        );
      },

      humanBytes(bytes, si) {
        let thresh = si ? 1000 : 1024;
        if (Math.abs(bytes) < thresh) {
          return bytes + 'B';
        }
        let units = si
          ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        do {
          bytes /= thresh;
          ++u;
        } while (Math.abs(bytes) >= thresh && u < units.length - 1);
        return bytes.toFixed(1) + '' + units[u];
      },

      forceRefresh() {
        this.refresh = false;
        this.$nextTick(() => {
          this.refresh = true;
        });
      },

      async getDirectories() {

        if (this.directories.length === 0) {

          let directories = [];

          this.$store.dispatch('loadDirectories', { path: '' })
            .then(response => {
              for (let i = 0; i < response.length; i++) {
                if (response[i].name !== '/') {
                  directories.push(response[i]);
                }
              }
            })
            .catch(error => console.warn(error));

          await pause(1000);

          for (let i = 0; i < directories.length; i++) {
            this.directories.push(directories[i]);
          }

          console.log("INITIAL: ",  JSON.stringify(this.directories, null, 1))
          this.showDirectories = !this.showDirectories
        }

        else {
          this.showDirectories = !this.showDirectories
        }

      },

      async loadList(table) {

        this.showDirectories = false;

        while (this.directories.length !== 0) {
          this.directories.pop();
        }

        this.itemscol = {};
        this.searchcol = {};

        this.active = [];
        this.open = [];
        this.searchField = '';

        this.$store.commit('setItem', {});
        if (!this.navigationRight.mini) {
          this.toogleRight();
        }
        this.$store.commit('setTable', table);
        this.$store.dispatch('loadItems')
          .then(() => {

            for (let i = 0; i < this.$store.state.headers.length; i++) {
              this.itemscol[this.$store.state.headers[i].value] = [];
            }

            this.forceRefresh();

            this.filter = {
              'table': this.$store.state.type,
              'columns': {}
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
        let sort = {
          'table': this.$store.state.type,
          'columns': {}
        };
        for (let i = 0; i < e.sortBy.length; i++) {
          if (e.sortDesc[i]) {
            sort['columns'][e.sortBy[i]] = 'DESC';
          } else {
            sort['columns'][e.sortBy[i]] = 'ASC';
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
        let type = this.$store.state.type;

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
          this.filter['table'] = type;
          this.filter['columns'][column] = [];
        } else {
          this.filter['table'] = type;
          this.filter['columns'][column] = [e[0] + slash];
        }

        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');
      },

      searchFilter() {
        this.filter['table'] = this.$store.state.type;

        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          if (key === 'origin') {
            column = 'origin.path';
          } else {
            column = key;
          }
          this.filter['columns'][column] = value;
        }

        this.$store.commit('setFilter', this.filter);
        this.$store.dispatch('loadItems');

      },

      clearFilter() {
        this.search = '';
        this.searchFilter();
      },

      async fetch(item) {

        let directories = [];

        this.$store.dispatch('loadDirectories', { path: item.path })
          .then(response => {
            if ((response.length === 1) && (response[0].name === '/')) {
              delete item.children;
            } else {
              for (let i = 0; i < response.length; i++) {
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
            vnode: null
          };
          this.$refs.treeView.nodes[child.path] = childNode;
        });

        for (let i = 0; i < directories.length; i++) {
          item.children.push(directories[i]);
        }

      },

    },

    mounted() {
      this.setBorderWidthLeft();
      this.setBorderWidthRight();
      this.setEventsAll();
    },

    destroyed() {
      this.$store.commit('setItem', {});
    }

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
</style>
