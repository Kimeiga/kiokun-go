<script lang="ts">
    export let data;
</script>

<div class="container">
    {#each data.entries as entry}
        <div class="entry">
            <h2>{entry.kanji.join('、 ')}</h2>
            <div class="reading">{entry.reading.join('、 ')}</div>
            {#each entry.sense as sense, i}
                <div class="sense">
                    {i + 1}. {sense.gloss.join('; ')}
                </div>
            {/each}
        </div>
    {/each}
</div>

<style>
    .container {
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
    }

    .entry {
        margin-bottom: 20px;
        padding: 15px;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .reading {
        color: #666;
        font-size: 0.9em;
    }

    .sense {
        margin-top: 10px;
    }
</style> 