<template>
    <div>
        <b-sidebar id="sidebar-1" bg-variant="dark" text-variant="light">
            <div>
                <b-button href="#" to="/" class="nav-button">Home</b-button>
                <b-button href="#" to="/routines" class="nav-button">Routines</b-button>
                <b-button href="#" to="/sensors" class="nav-button">Sensors</b-button>
                <b-button href="#" to="/recipes" class="nav-button">Recipes</b-button>
                <b-button href="#" to="/inventory" class="nav-button">Inventory</b-button>
                <b-button href="#" to="/cameras" class="nav-button">Cameras</b-button>
            </div>
        </b-sidebar>
    </div>
</template>

<script>
  export default {
    name: "SideBar"
  }
</script>

<style scoped>
    .nav-button {
        width: 320px;
        margin-top: 10px;
        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);
    }
</style>
