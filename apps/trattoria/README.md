# Trattoria

A modern restaurant website built with Next.js, featuring a responsive design showcasing menu items, about section, and contact information. This project utilizes Contentful as a headless CMS to manage restaurant menu items and marketing assets.

## Features

- **Responsive Design** - Optimized for all device sizes
- **Dynamic Menu** - Menu items pulled from Contentful CMS
- **Component-Based Architecture** - Modular sections for easy maintenance
- **Marketing Assets** - Images and content managed through Contentful

## Tech Stack

- [Next.js 15](https://nextjs.org) with App Router
- [React 19](https://react.dev/)
- [Tailwind CSS 4](https://tailwindcss.com/)
- [Contentful](https://www.contentful.com/) for content management
- [Turbopack](https://turbo.build/) for fast development experience

## Getting Started

### Prerequisites

- Node.js 18.17.0 or later
- pnpm (recommended for this monorepo)
- Contentful account with appropriate space setup

### Environment Variables

Create a `.env.local` file in the project root with the following variables:

```
CONTENTFUL_SPACE_ID=your_space_id
CONTENTFUL_ACCESS_TOKEN=your_access_token
```

### Installation

This project is part of a monorepo managed with pnpm workspaces. From the monorepo root:

```bash
# Install dependencies
pnpm install

# Start the development server
pnpm --filter trattoria dev
```

Or, from the app directory:

```bash
cd apps/trattoria
pnpm dev
```

The application will be available at [http://localhost:3000](http://localhost:3000).

## Project Structure

- `app/` - Next.js app directory
  - `components/` - React components organized by section
  - `lib/` - Utility functions including Contentful integration
  - `page.tsx` - Main page component

## Contentful Setup

This project uses Contentful to manage:
- Menu items (name, description, price, category, image)
- Marketing assets (hero images, about section content)

## Scripts

- `pnpm dev` - Start the development server with Turbopack
- `pnpm build` - Build the application for production
- `pnpm start` - Start the production server
- `pnpm lint` - Run ESLint with zero tolerance for warnings
- `pnpm check-types` - Run TypeScript type checking

## Deployment

This application can be deployed on [Vercel](https://vercel.com) or any other Next.js-compatible hosting service.

## Contributing

Please see the monorepo README for contribution guidelines.