<script>
  import 'bootstrap/dist/css/bootstrap.min.css'
  import 'bootstrap/dist/js/bootstrap.bundle'

  export let subscription, run_vpn_function, delete_subs_function;
</script>

<main>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">{subscription["display_name"]}</h5>
      <p class="card-text">{subscription["url"]}</p>
      <button class="btn btn-primary" on:click={()=>{run_vpn_function(subscription)}}>Connect</button>
      <button class="btn btn-danger" on:click={()=>{delete_subs_function(subscription)}}>Delete</button>
    </div>
  </div>
</main>

<style>
</style>
