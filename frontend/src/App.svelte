<script>
  import 'bootstrap/dist/css/bootstrap.min.css'
  import logo from './assets/images/ecorp-logo.png'

  import SubscriptionCard from './SubscriptionCard.svelte'
  import bootstrapBundle from 'bootstrap/dist/js/bootstrap.bundle'
  import {EventsOn} from "../wailsjs/runtime/runtime"

  import {ParseConfig, StartVPN, WriteConfig} from '../wailsjs/go/main/App'
  import {main} from "../wailsjs/go/models"

  let resultText = "Fill to add another subscription!"
  let newUrl
  let configName
  let configs = []
  let resultSuccess
  let resultFailure
  let toastElement

  document.addEventListener("DOMContentLoaded", async () => {
    let tmp = await ParseConfig()
    configs = tmp["subscriptions"]
  });

  $: if(resultFailure != null) {
    console.log(resultFailure)
    toastElement = bootstrapBundle.Toast.getOrCreateInstance(document.getElementById("toast-failure"))
    toastElement.show();
  }

  $: if(resultSuccess != null) {
    console.log(resultSuccess)
    toastElement = bootstrapBundle.Toast.getOrCreateInstance(document.getElementById("toast-success"))
    toastElement.show();
  }

  async function RunVPN(cfg) {
    StartVPN(cfg).then(
      () => {
        resultSuccess = "Connected"
      },
      (reason) => {
        resultFailure = "Could not connect. Check installation"
        console.log(reason)
      }
    )
  }

  EventsOn("VPNLog", (log) => {console.log(log)})

  async function AddURL() {
    let url = newUrl
    let res = fetch(url, {method: "GET"})
    res.then(async (result) => {
      let cfg = await result.json()
      let curr_date = new Date()
      let subscription = {
        last_update: curr_date.getTime(),
        display_name: configName,
        url: url,
        config: cfg
      }
      configs.push(subscription)
      configs = configs
      let ok = await WriteConfig(new main.AppConfig({subscriptions: configs}))
      resultText = ok ? "Added successfully" : "Writing the new config file failed :("
      if(ok) {
        resultSuccess = resultText
      } else {
        resultFailure = resultText
      }
    }, (reason) => {
      console.log("URL Fetch Failed:", reason)
      resultText = "Failed to add subscription URL: " + reason
      resultFailure = resultText
    })
  }

  async function DeleteSubscription(subscription) {
    let target = JSON.stringify(subscription);
    let newConfig = [];
    configs.forEach(element => {
      if(target !== JSON.stringify(element)) {
        newConfig.push(element)
      }
    })
    configs = newConfig
    let ok = await WriteConfig(new main.AppConfig({subscriptions: configs}))
    resultText = ok ? "Deleted successfully" : "Writing the new config file failed :("
    if(ok) {
      resultSuccess = resultText
    } else {
      resultFailure = resultText
    }
  }
</script>

<main>
  <img alt="Wails logo" id="logo" src="{logo}">
  <div class="card my-4" style="min-height: 12rem">
    <div class="card-header card-title text-center">Subscriptions</div>
    <div class="card-body">
      <div class="row">
        {#each configs as subs, _i}
          <div class="col-3">
            <SubscriptionCard subscription={subs} run_vpn_function={RunVPN} delete_subs_function={DeleteSubscription}/>
          </div>
        {/each}
      </div>
    </div>
  </div>
  <div class="card">
    <!-- New config card -->
    <div class="card-header card-title text-center">Add a new subscription</div>
    <div class="card-body align-items-center">
      <form class="my-3">
        <div class="row align-items-center">
          <div class="col-7">
            <input autocomplete="on" placeholder="URL" bind:value={newUrl} class="form-control" id="cfgurl" type="url"/>
          </div>
          <div class="col-3">
            <input autocomplete="on" placeholder="Name" bind:value={configName} class="form-control" id="cfgname" type="text"/>
          </div>
          <div class="col-auto">
            <button class="btn btn-primary" type="button" on:click={AddURL}>Add URL</button>
          </div>
        </div>
      </form>
    </div>
  </div>
  <div class="position-relative bottom-0 end-0">
    <div class="toast-container p-3 bottom-0 end-0">
      <div class="toast align-items-center text-bg-primary border-0" id="toast-success" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="d-flex">
          <div class="toast-body">
            {resultSuccess}
          </div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
      <div class="toast align-items-center text-bg-danger border-0" id="toast-failure" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="d-flex">
          <div class="toast-body">
            {resultFailure}
          </div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
    </div>
  </div>
</main>

<style>

  #logo {
    display: block;
    width: 20%;
    height: 20%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

</style>
