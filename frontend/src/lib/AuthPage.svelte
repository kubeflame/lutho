<script lang="ts">
  import HeaderElement from "./HeaderElement.svelte";
  import { type AuthResponse } from "./types";
  import { apiURL, KubeAuthType, routeString, transitionEffects } from "./util";
  import { push, querystring } from "svelte-spa-router";

  let errorMessage: string = "";
  let referrerQueryParam =
    (new URLSearchParams($querystring)
      .get("referrer")
      ?.toLowerCase() as string) || routeString.home;

  let authType: keyof typeof KubeAuthType = "accessToken";
  let accessToken: string = "";
  let kubeconfigPath: string = "";
  let kubeconfigRaw: string = "";
  let masterURL: string = "";

  $: authRequest = {
    type: authType,
    token: accessToken,
    kubeconfigPath: kubeconfigPath,
    kubeconfigRaw: kubeconfigRaw,
    masterURL: masterURL,
  };

  async function onClickAuthButton() {
    await fetch(apiURL.auth, {
      method: "POST",
      body: JSON.stringify(authRequest),
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((result: AuthResponse) => {
        errorMessage = result.error;
        if (!errorMessage) {
          push(referrerQueryParam);
        }
      })
      .catch((error) => {
        errorMessage = error;
      });
  }
</script>

<HeaderElement
  headerContentClassName={"shadow-md rounded-b-2xl bg-base-100"}
  logoClassName={""}
/>

<div
  class="auth bg-base-100 fixed top-10 m-4 flex h-[calc(100%-72px)] w-[calc(100%-32px)]
    flex-col rounded-md shadow-md {transitionEffects}"
>
  <div class="join flex h-6 w-full items-center">
    <label
      for="authType-select"
      class="join-item bg-base-300 h-full p-0.5 pl-2 pr-2 text-sm font-normal"
    >
      Auth Type
    </label>
    <select
      id="authType-select"
      bind:value={authType}
      class="join-item select select-primary select-xs bg-base-100 mt-1 flex h-full rounded-full text-center
        text-sm font-normal drop-shadow-sm focus:outline-0"
    >
      {#each Object.keys(KubeAuthType) as at}
        <option value={at}>{at}</option>
      {/each}
    </select>
  </div>
  {#if authType === KubeAuthType.accessToken}
    <div class="join h-6 w-full items-center">
      <label
        for="masterURL-input"
        class="join-item bg-base-300 h-full p-0.5 pl-2 pr-2 text-sm font-normal"
      >
        Master URL
      </label>
      <input
        id="masterURL-input"
        bind:value={masterURL}
        type="text"
        class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
          flex grow text-sm focus:outline-0"
      />
    </div>
    <div class="join h-6 w-full items-center">
      <label
        for="accessToken-input"
        class="join-item bg-base-300 h-full p-0.5 pl-2 pr-2 text-sm font-normal"
      >
        AccessToken
      </label>
      <input
        id="accessToken-input"
        bind:value={accessToken}
        type="password"
        class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
    flex grow text-sm focus:outline-0"
      />
    </div>
  {:else if authType === KubeAuthType.kubeconfigPath}
    <div class="join h-6 w-full items-center">
      <label
        for="kubeconfigPath-input"
        class="join-item bg-base-300 h-full p-0.5 pl-2 pr-2 text-sm font-normal"
      >
        Kubeconfig Path
      </label>
      <input
        id="kubeconfigPath-input"
        bind:value={kubeconfigPath}
        type="text"
        class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
flex grow text-sm focus:outline-0"
      />
    </div>
  {:else if authType === KubeAuthType.kubeconfigRaw}
    <textarea
      bind:value={kubeconfigRaw}
      class="bg-base-200 outline-base-100 focus:outline-primary h-36 min-w-full resize-none
        overflow-y-scroll rounded-lg p-2 shadow-md outline outline-1"
    />
  {/if}
  <div
    class="auth-content hero-content relative h-full flex-col overflow-y-scroll lg:flex-row-reverse"
  >
    <div class="text-center lg:text-left">
      <h1 class="text-5xl font-bold">Login now!</h1>
      <p class="py-6">
        {errorMessage}
      </p>
      <button on:click={onClickAuthButton} class="btn btn-primary">
        Login
      </button>
    </div>
  </div>
</div>
