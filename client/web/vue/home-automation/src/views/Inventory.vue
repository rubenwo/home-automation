<template>
    <div style="background-color: rgba(255, 255, 255, 0.7);">
        <div v-if="!this.inventory.length <= 0">
            <div v-bind:key="item.id" v-for="item in this.inventory">
                <b-row>
                    <h3>{{item.name}}</h3>
                    <p>{{item.description}}</p>
                    <p>{{item.amount}}</p>
                    <b-button @click="removeItem(item.id)">X</b-button>
                </b-row>
            </div>
        </div>
        <div v-else>Nothing here yet</div>


    </div>
</template>

<script>
  import InventoryService from '../services/inventory.service';

  export default {
    name: "Inventory",
    data() {
      return {
        inventory: []
      }
    },
    methods: {
      async removeItem(id) {
        console.log("removing item with id: " + id);
        const res = await InventoryService.deleteInventoryItem(id);
        console.log(res);
        this.inventory = res;
      }
    },
    async mounted() {
      const res = await InventoryService.fetchInventory();
      this.inventory = res;
    },
  }
</script>

<style scoped>

</style>
