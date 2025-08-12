# 🗺️ Port Creation Guide

A port is applying OneDark's palette for an application. Essentially, it's a theme for that software's UI elements.

## Table of Contents
- [🗺️ Port Creation Guide](#️-port-creation-guide)
  - [Table of Contents](#table-of-contents)
  - [Creating the Repository](#creating-the-repository)
  - [Configuring the Repository](#configuring-the-repository)
  - [Submit Your Port](#submit-your-port)
  - [Licensing](#licensing)
  - [Tools](#tools)

## Creating the Repository

Ports can be created using the OneDark repository template found [here](https://github.com/onedarktheme/template).

> [!IMPORTANT]
> Make sure that the name of your repository is the name of the application you are porting to OneDark in `lower-kebab-case`

Steps:
1. Clone the template repository
```bash
git clone https://github.com/catppuccin/template.git <name-of-the-application-being-ported>
```

2. Navigate to the root directory of your new repository
```bash
cd <name-of-the-application-being-ported
```

3. Re-initialize the repository to remove the template authors from the contributors of your port
```bash
rm -rf ./.git/
git init
```

## Configuring the Repository

1. Make sure that the default branch is `master` for consistency across all OneDark ports
2. Put all images in the `assets/` directory
   - Use `.webp` as the image format when possible since it offers better compression
3. Set repository description to:
```text
<emoji> dark terminally stylish theme for <Application Name>
```
    - `<emoji>` is an emoji that best represents the app
    - `<Application Name>` is the properly capitalized name of the application you are porting
4. Add `onedark`, `theme`, and `name-of-the-application-being-ported` to the repository topics
5. Update all shields.io badges to point to your repository instead of the template
6. Add preview screenshots for all 4 palettes (use [relative links](https://github.blog/news-insights/product-news/relative-links-in-markup-files/))
7. Add yourself to the `🙏 Acknowledgements` section

> [!CAUTION]
> Do **not** add organization-wide files such as `CODE_OF_CONDUCT.md` because they are maintained in [this repository](https://github.com/onedarktheme/.github)

## Submit Your Port

Open an issue [here](https://github.com/onedarktheme/onedark/issues/new?assignees=&template=port-submission.yml&title=Name+of+the+application+being+ported) to get the port added to the organization.

## Licensing

All ports use the MIT license.

## Tools

You can look for useful tools (or contribute new ones with a PR!) in the [OneDark toolbox](https://github.com/onedarktheme/toolbox) to help you create your port.