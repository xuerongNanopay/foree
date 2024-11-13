<script lang="ts">

</script>

<ul class="accordion">
    <li>
        <input type="radio" name="accordion" id="first" checked>
        <label for="first">Menu</label>
        <ul class="content">
            <li>Dashboard</li>
            <li>Send</li>
            <li>Send</li>
            <li>Send</li>
            <li>Send</li>
        </ul>
    </li>
    <br/>
    <li>
        <input type="radio" name="accordion" id="second">
        <label for="second">Profile</label>
        <ul class="content">
            <li>me</li>
            <li>About</li>
            <li>About</li>
            <li>About</li>
            <li>About</li>
            <li>About</li>
        </ul>
    </li>
</ul>

<style>
    .accordion {
        background: salmon;
    }

    .accordion input[type="radio"] {
        display: none;
    }

    .accordion .content {
        max-height: 0;
        overflow: hidden;
        transition: max-height 1s ease-in-out;
    }

    .accordion input[type="radio"]:checked + label + .content {
        max-height: 600px;
    }

    .accordion input[type="radio"] + label::before {
        content: "+";
        margin-right: 10px;
        font-size: 24px;
        font-weight: 600;
    }

    .accordion input[type="radio"]:checked + label::before {
        content: "-";
    }
</style>