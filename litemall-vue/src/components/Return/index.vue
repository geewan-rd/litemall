<template>
	<div class="button-return" v-if="show" v-on:click="click">
		<van-icon name="arrow-left" class="icon" />
	</div>
</template>

<script>
export default {
  name: 'button-return',
  computed: {
    show: function () {
      const paths = ['/', '/items', '/order', '/user']
      return !paths.includes(this.$route.path)
    }
  },
  methods: {
    click() {
      this.$router.go('-1')
    }
  }
};
</script>


<style lang="scss" scoped>
.button-return {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  bottom: 70px;
  right: 20px;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: rgba($color: #000000, $alpha: 0.2);
  z-index: 1000;
}
.icon {
  color: #fff;
}
</style>
