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
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <v-card class="mt-10" tile>
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
        <template v-slot:item.accessed="{ item }">
          <div>{{ toLocal(item.accessed) }}</div>
        </template>
        <template v-slot:item.modified="{ item }">
          <div>{{ toLocal(item.modified) }}</div>
        </template>
        <template v-slot:item.created="{ item }">
          <div>{{ toLocal(item.created) }}</div>
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

  export default {
    name: 'items',
    components: {
      JsonToHtml,
      item
    },
    data: function () {
      return {
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
      }
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
      loadList(table) {
        this.$store.commit('setItem', {});
        if (!this.navigationRight.mini) {
          this.toogleRight();
        }
        this.$store.commit('setTable', table);
        this.$store.dispatch('loadItems');

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

  .fade-enter-active, .fade-leave-active {
    transition: opacity .2s;
  }

  .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */
  {
    opacity: 0;
  }
</style>
