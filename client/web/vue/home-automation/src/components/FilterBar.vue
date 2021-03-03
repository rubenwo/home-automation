<template>
  <div class="filter-bar ui basic segment grid">
    <div class="ui form">
      <div class="inline field">
        <label>Search for:</label>
        <b-input
          class="sm-1"
          type="text"
          v-model="filterText"
          @keyup.enter="doFilter"
          placeholder="Search..."
        />
        <b-button variant="success" @click="doFilter">Go</b-button>
        <b-button @click="resetFilter">Reset</b-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      filterText: ""
    };
  },
  methods: {
    doFilter() {
      this.$events.fire("filter-set", this.filterText);
    },
    resetFilter() {
      this.filterText = ""; // clear the text in text input
      this.$events.fire("filter-reset");
    }
  }
};
</script>
