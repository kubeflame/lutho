<script lang="ts">
</script>

<div class="h-full w-full">
  <div
    class="spinner newton flex h-2/4 place-content-center items-end drop-shadow-md"
    id="newton"
  >
    <div class="ball drop-shadow-sm" />
    <div class="ball drop-shadow-sm" />
    <div class="ball drop-shadow-sm" />
    <div class="ball drop-shadow-sm" />
    <div class="ball drop-shadow-sm" />
  </div>

  <div
    class="to-primary text-primary flex place-content-center items-center space-x-2 bg-gradient-to-b from-transparent bg-clip-text pt-10 font-mono text-3xl text-transparent drop-shadow-md"
  >
    <span>L</span>
    <span>O</span>
    <span>A</span>
    <span>D</span>
    <span>I</span>
    <span>N</span>
    <span>G</span>
    <span>.</span>
    <span>.</span>
    <span>.</span>
  </div>
</div>

<style>
  @-webkit-keyframes newtonBallFirst {
    0% {
      -webkit-transform: rotate(70deg);
    }
    50% {
      -webkit-transform: rotate(0deg);
    }
    100% {
      -webkit-transform: rotate(0deg);
    }
  }
  @-webkit-keyframes newtonBallMiddle {
    0% {
      -webkit-transform: rotate(0deg);
    }
    50% {
      -webkit-transform: rotate(0deg);
    }
    51% {
      -webkit-transform: rotate(-0.5deg);
    }
    52% {
      -webkit-transform: rotate(0.5deg);
    }
    53% {
      -webkit-transform: rotate(0deg);
    }
    100% {
      -webkit-transform: rotate(0deg);
    }
  }
  @-webkit-keyframes newtonBallLast {
    0% {
      -webkit-transform: rotate(0deg);
    }
    50% {
      -webkit-transform: rotate(0deg);
    }
    100% {
      -webkit-transform: rotate(-70deg);
    }
  }
  @keyframes newtonBallFirst {
    0% {
      transform: rotate(70deg);
    }
    50% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(0deg);
    }
  }
  @keyframes newtonBallMiddle {
    0% {
      transform: rotate(0deg);
    }
    50% {
      transform: rotate(0deg);
    }
    51% {
      transform: rotate(-0.5deg);
    }
    52% {
      transform: rotate(0.5deg);
    }
    53% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(0deg);
    }
  }
  @keyframes newtonBallLast {
    0% {
      transform: rotate(0deg);
    }
    50% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(-70deg);
    }
  }

  .spinner.newton .ball {
    position: relative;
    display: inline-block;
    width: 1em;
    height: 6em;
    -webkit-transform-origin: 50% 0%;
    transform-origin: 50% 0%;
  }

  .spinner.newton .ball::after {
    content: "";
    position: absolute;
    left: 0;
    bottom: 0;
    width: 1em;
    height: 1em;
    border-radius: 100%;
    /* background: radial-gradient(
            circle at 40% 40%,
            rgb(79, 77, 77) 39%,
            rgb(55, 55, 55) 42%
        ); */
    /* background: linear-gradient(oklch(var(--p)), transparent); */
    background: radial-gradient(transparent, oklch(var(--p)));
  }

  .spinner.newton .ball::before {
    content: "";
    position: absolute;
    width: 0.08em;
    margin-left: -0.04em;
    top: 0;
    left: 50%;
    bottom: 1em;
    background: linear-gradient(transparent, oklch(var(--p)));
  }

  .spinner.newton .ball {
    -webkit-animation: newtonBallMiddle 1s infinite alternate;
    animation: newtonBallMiddle 1s infinite alternate;
  }

  .spinner.newton .ball:first-child {
    -webkit-animation: newtonBallFirst 1s ease-in infinite alternate;
    animation: newtonBallFirst 1s ease-in infinite alternate;
  }
  .spinner.newton .ball:first-child::after {
    -webkit-animation: newtonBallLast 1s ease-out infinite alternate-reverse;
    animation: newtonBallLast 1s ease-out infinite alternate-reverse;
  }

  .spinner.newton .ball:last-child {
    -webkit-animation: newtonBallLast 1s ease-out infinite alternate;
    animation: newtonBallLast 1s ease-out infinite alternate;
  }
  .spinner.newton .ball:last-child::after {
    -webkit-animation: newtonBallFirst 1s ease-in infinite alternate-reverse;
    animation: newtonBallFirst 1s ease-in infinite alternate-reverse;
  }
</style>
