<script lang="ts">
  import type { Person } from "$lib/types";
  import type { PageData } from "./$types";
  import autoAnimate from "@formkit/auto-animate"

  export let data: PageData;

  const handleAddPerson = async () => {
    const person : Person = await fetch(`/persons`, {method: "POST"}).then(res => res.json());
    data.persons = [...data.persons, person];
  };
</script>

<h1>👤 People list</h1>
<div class="row flex-spaces child-borders">
  <button on:click={handleAddPerson}>Add random person</button>
</div>
<ul use:autoAnimate={{ duration: 150 }}>
  {#each data.persons as person}
    <li>
      {person.last_name} 
      <ul>
        <li>📱 {person.phone_number}</li>
        <li>📍 {person.location}</li>
      </ul>
    </li>
    <br>
  {/each}
</ul>
