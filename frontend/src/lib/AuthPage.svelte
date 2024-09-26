<script lang="ts">
  import { onMount } from "svelte";
  import BottomAlert from "./BottomAlert.svelte";
  import HeaderElement from "./HeaderElement.svelte";
  import { type Alert, type AuthResponse } from "./types";
  import { apiURL, KubeAuthType, routeString, transitionEffects } from "./util";
  import { push, querystring } from "svelte-spa-router";
  import { authState, kubeHost } from "./stores";

  let referrerQueryParam =
    (new URLSearchParams($querystring)
      .get("referrer")
      ?.toLowerCase() as string) || routeString.home;

  let authType: keyof typeof KubeAuthType = "accessToken";
  let accessToken: string = "";
  let kubeconfigPath: string = "";
  let kubeconfigRaw: string = "";
  let masterURL: string = "";

  $: alert = { message: null, type: null } as Alert;
  $: authRequest = {
    type: authType,
    token: accessToken,
    kubeconfigPath: kubeconfigPath,
    kubeconfigRaw: kubeconfigRaw,
    masterURL: masterURL,
  };

  onMount(async () => {
    if ($authState && $kubeHost) push(referrerQueryParam);
  });

  async function onClickAuthButton() {
    await fetch(`${location.pathname}${apiURL.auth}`, {
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
        alert = { message: result.error, type: "error" };
        if (!alert.message) {
          $authState = result.state;
          $kubeHost = result.kubeHost;
          push(referrerQueryParam);
        }
      })
      .catch((error) => {
        alert = { message: error, type: "error" };
      });
  }

  const authTypeData = [
    {
      name: "Access Token",
      type: KubeAuthType.accessToken,
    },
    {
      name: "Kubeconfig Path",
      type: KubeAuthType.kubeconfigPath,
    },
    {
      name: "Kubeconfig Raw",
      type: KubeAuthType.kubeconfigRaw,
    },
  ];
</script>

<HeaderElement
  headerContentClassName={"shadow-md rounded-b-2xl bg-base-100"}
  logoClassName={""}
/>

<div
  class="auth fixed top-10 m-4 flex h-[calc(100%-72px)] w-[calc(100%-32px)] content-center
    justify-center rounded-md bg-base-100 shadow-md {transitionEffects}"
>
  <div class="h-1/2 w-1/2 space-y-2 place-self-center">
    <div class="join flex h-6 w-full place-items-center items-center">
      <label
        for="authType-select"
        class="join-item flex h-full bg-base-200 pl-2 pr-2 text-sm font-normal"
      >
        Auth Type
      </label>
      <select
        id="authType-select"
        bind:value={authType}
        class="join-item select select-xs flex rounded-lg border-base-300 bg-base-100
          text-center text-sm font-normal drop-shadow-sm focus:border-primary/60 focus:outline-0"
        on:change={() => {
          alert = { message: null, type: null };
        }}
      >
        {#each authTypeData as at}
          <option value={at.type}>{at.name}</option>
        {/each}
      </select>
    </div>
    {#if authType === KubeAuthType.accessToken}
      <div class="join h-6 w-full items-center">
        <label
          for="masterURL-input"
          class="join-item h-full bg-base-200 p-0.5 pl-2 pr-2 text-sm font-normal"
        >
          Master URL
        </label>
        <input
          id="masterURL-input"
          bind:value={masterURL}
          type="text"
          class="input input-xs join-item input-bordered flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>
      <div class="join h-6 w-full items-center">
        <label
          for="accessToken-input"
          class="join-item h-full bg-base-200 p-0.5 pl-2 pr-2 text-sm font-normal"
        >
          AccessToken
        </label>
        <input
          id="accessToken-input"
          bind:value={accessToken}
          type="password"
          class="input input-xs join-item input-bordered flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>
    {:else if authType === KubeAuthType.kubeconfigPath}
      <div class="join h-6 w-full items-center">
        <label
          for="kubeconfigPath-input"
          class="join-item h-full bg-base-200 p-0.5 pl-2 pr-2 text-sm font-normal"
        >
          Kubeconfig Path
        </label>
        <input
          id="kubeconfigPath-input"
          bind:value={kubeconfigPath}
          type="text"
          class="input input-xs join-item input-bordered flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>
    {:else if authType === KubeAuthType.kubeconfigRaw}
      <textarea
        bind:value={kubeconfigRaw}
        class="h-full min-w-full resize-none overflow-y-scroll rounded-lg bg-base-200/20
          p-2 shadow-md outline outline-1 outline-base-200 focus:outline-primary"
      />
    {/if}

    <div class="flex w-full content-center items-center justify-center pt-5">
      <button
        on:click={onClickAuthButton}
        class="btn btn-ghost btn-primary btn-sm bg-base-200 hover:bg-primary hover:drop-shadow-md"
      >
        Login
      </button>
      <BottomAlert {alert} />
    </div>
  </div>
</div>
