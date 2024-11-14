<script lang="ts">
    import UpArrow from "$lib/assets/icons/up_arrow.png"
</script>

<ul class="accordion">
    <li>
        <input type="checkbox" name="first" id="first" checked>
        <label for="first">
            <p>Menu</p>
            <img src={UpArrow} alt="">
        </label>
        <ul class="content">
            <li>Dashboard</li>
            <li>Send</li>
            <li>Send</li>
            <li>Send</li>
            <li>Send</li>
        </ul>
    </li>
    <li>
        <input type="checkbox" name="second" id="second">
        <label for="second">
            <p>Menu</p>
            <img src={UpArrow} alt="">
        </label>
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
        padding: 0 1rem;
    }

    .accordion label {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: .5rem 0 .3rem;
        border-bottom: 1px solid var(--slate-200);
        &:hover {
            background-color: var(--foree-bg-4);
        }

        & p {
            font-size: large;
            font-weight: 500;
        }
        & img {
            height: 17px;
            width: 17px;
            transition: all .5s ease-in-out;
        }
    }

    .accordion > li > ul > li{
        padding: .5rem 1rem;
        border-radius: 7px;
        &:hover {
            background-color: var(--foree-bg-4);
        }
    }

    .accordion input[type="checkbox"] {
        display: none;
    }

    .accordion .content {
        max-height: 0;
        overflow: hidden;
        transition: max-height .5s ease-in-out;
    }

    .accordion input[type="checkbox"]:checked {
        & ~ .content {
            max-height: 300px;
        }

        & ~ label > img {
            transform: rotate(180deg);
        }
    }
</style>