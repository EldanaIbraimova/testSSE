<script>
import { onMount } from "svelte"
let time = ""
onMount( ( => {
    const evtSrc = new EventSource("http://localhost:3500/event")
}))
</script>

<main>

</main>