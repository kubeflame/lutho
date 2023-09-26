<script lang="ts">
  import HeaderElement from "./HeaderElement.svelte";
</script>

<HeaderElement
  headerContentClassName={"shadow-md rounded-b-2xl bg-base-100"}
  logoClassName={""}
/>

<div class="auth hero fixed top-10 flex h-screen flex-col bg-base-100">
  <div
    class="auth-content overflow-y-scroll; hero-content relative h-full flex-col lg:flex-row-reverse"
  >
    <div class="text-center lg:text-left">
      <h1 class="text-5xl font-bold">Login now!</h1>
      <p class="py-6">
        Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda
        excepturi exercitationem quasi. In deleniti eaque aut repudiandae et a
        id nisi.
      </p>
    </div>
    <div class="card w-full max-w-sm flex-shrink-0 bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Email</span>
            <input
              type="text"
              placeholder="email"
              class="input input-bordered"
            />
          </label>
        </div>
        <div class="form-control">
          <label class="label">
            <span class="label-text">Password</span>
            <input
              type="text"
              placeholder="password"
              class="input input-bordered"
            />
          </label>
        </div>
        <div class="form-control mt-6">
          <button class="btn btn-primary">Login</button>
        </div>
      </div>
    </div>
  </div>
</div>
