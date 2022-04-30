import Link from 'next/link';
import { Menu } from '@headlessui/react';

export default function Sidebar() {
  return (
    <aside className="h-screen sticky top-0 flex flex-col items-center w-60 text-gray-400 bg-gray-900">
      <a className="flex items-center text w-full pt-8 px-4" href="#">
        <svg
          viewBox="0 0 1024 1024"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          className="w-10 h-10 fill-current"
        >
          {/* <rect width="1024" height="1024" rx="221.867" fill="white" /> */}
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M338.091 302.873L338.109 302.856L338.127 302.839C357.817 283.839 381.03 269.324 404.994 260.768L405.018 260.76L405.042 260.751C426.583 253.003 450.005 249.564 472.609 250.804C491.395 251.835 509.477 255.778 526.373 262.494L526.434 262.518L526.495 262.542C539.783 267.73 552.046 274.295 562.909 281.983L563.015 282.058L563.123 282.132C564.473 283.063 565.803 284.018 567.117 284.998C565.831 285.608 564.524 286.233 563.197 286.871L563.025 286.954L562.854 287.04C554.269 291.36 542.16 297.642 529.564 305.696L529.482 305.749L529.401 305.802C515.813 314.664 504.316 323.553 494.655 332.952C483.752 343.469 474.437 355.121 467.047 367.705L467.032 367.729L467.018 367.754C460.582 378.78 455.33 390.768 450.471 405.036L450.467 405.049L450.462 405.061C446.01 418.188 442.207 431.623 438.842 446.062C435.77 459.131 432.767 474.234 429.092 494.824L428.953 495.603L428.861 496.492C428.735 497.505 428.555 498.537 428.194 500.607C427.803 502.794 427.382 505.156 426.901 508.343L426.892 508.401L426.28 512.321L426.274 512.352L425.993 514.024C425.813 515.089 425.569 516.532 425.302 518.089C424.749 521.318 424.159 524.679 423.83 526.338L423.774 526.619L423.726 526.901C421.72 538.63 418.381 551.792 413.443 566.936C404.038 595.068 390.145 620.752 373.501 641.31L373.472 641.346L373.443 641.383C356.357 662.674 334.147 680.624 311.343 691.83L311.275 691.864L311.206 691.898C290.614 702.201 267.621 708.088 244.922 708.896L244.899 708.897L244.877 708.898C226.133 709.604 207.763 707.249 190.267 701.93C176.338 697.674 163.616 692 152.38 685.152L152.339 685.128L152.299 685.103C150.195 683.834 148.141 682.525 146.147 681.186C147.965 680.056 149.844 678.872 151.785 677.628L151.852 677.586L151.918 677.542C160.87 671.714 172.108 664.161 183.111 655.021C195.242 645.002 205.806 634.632 214.353 624.083C223.914 612.49 231.583 599.919 237.122 586.633C241.904 575.351 245.466 563.666 248.76 549.204L248.767 549.175L248.774 549.146C251.525 536.909 253.85 523.609 255.889 508.603L255.899 508.542C256.52 505.398 256.802 502.347 256.959 500.641C256.974 500.478 256.988 500.327 257.001 500.189C257.14 498.754 257.233 497.832 257.334 497.045C258.049 493.175 258.411 489.185 258.665 486.379L258.665 486.378L258.672 486.304L258.679 486.226L258.951 483.342L258.965 483.161C259.046 482.117 259.211 480.778 259.551 478.607L259.591 478.353L260.936 468.197L260.968 467.895C261.291 464.782 262.482 457.242 264.061 449.006L264.067 448.978L264.072 448.95C269.356 420.988 277.977 394.666 289.596 370.751L289.603 370.736L289.611 370.721C302.281 344.54 319.134 321.082 338.091 302.873ZM593.447 766.358L593.458 766.36L593.469 766.362C611.351 769.923 629.381 770.346 645.472 767.8L645.481 767.799C661.925 765.204 678.192 758.986 691 750.588L691.037 750.564C702.779 742.899 713.119 732.724 720.822 721.28C727.311 711.631 732.057 701.337 734.976 690.691L734.999 690.607L735.022 690.524C736.492 685.328 737.58 680.012 738.242 674.79C733.113 677.125 728.186 679.178 723.317 680.976L723.227 681.009L723.137 681.042C712.943 684.702 703.287 687.344 693.996 688.766C683.52 690.413 673.05 690.696 662.841 689.504L662.815 689.501L662.789 689.498C653.876 688.436 645.05 686.306 635.381 683.096C626.531 680.176 617.784 676.77 608.671 672.705C601.669 669.607 594.368 666.079 587.841 662.828L587.751 662.998L573.672 655.568C573.086 655.264 572.464 654.937 571.787 654.573L568.216 652.688C565.841 651.434 562.054 649.435 560.311 648.586L560.073 648.47L559.838 648.348C553.547 645.09 545.717 641.778 536.108 638.448C518.392 632.486 500.371 629.76 484.013 630.302L483.962 630.303L483.91 630.305C467.113 630.781 450.014 635.057 436.123 642.053L436.056 642.086L435.989 642.119C420.145 649.96 409.49 660.596 403.43 668.046L403.41 668.07L403.39 668.094C396.159 676.931 390.552 686.726 386.692 697.187C384.529 703.112 382.925 709.042 381.877 714.853C386.66 713.241 392.912 711.29 399.678 709.701C410.049 707.244 420.123 705.731 429.569 705.404C440.096 704.973 450.387 705.982 460.275 708.53C470.613 711.151 479.279 714.816 485.713 717.693L485.786 717.726L485.859 717.76C493.595 721.289 501.599 725.409 510.165 730.244C512.171 731.231 514.016 732.356 514.995 732.953L514.995 732.953L515.202 733.079L516.913 734.106C519.425 735.397 521.783 736.862 523.439 737.9L523.475 737.923L524.837 738.77L524.846 738.775C525.47 739.162 526.173 739.572 527.492 740.295L527.627 740.369L533.706 743.816L533.954 743.966C535.439 744.862 539.94 747.254 544.321 749.415L544.351 749.429L544.381 749.444C560.425 757.419 576.977 763.096 593.447 766.358ZM578.799 633.894L578.799 633.894L578.798 633.897L578.799 633.894ZM859.329 562.223L859.32 562.204C848.824 541.611 833.4 521.493 814.639 504.241L814.628 504.231L814.618 504.221C797.493 488.424 777.596 474.566 755.484 463.127L755.462 463.115L755.439 463.104C748.932 459.715 742.912 456.829 740.383 455.78L740.137 455.678L731.961 452.009L731.759 451.912C730.024 451.083 728.944 450.61 728.089 450.282L727.941 450.225L725.592 449.276L725.528 449.251L725.467 449.227L725.467 449.226C723.179 448.313 719.924 447.013 716.86 445.443C716.225 445.161 715.476 444.853 714.308 444.378L714.147 444.313L713.939 444.23C712.549 443.672 710.062 442.673 707.583 441.364L707.534 441.341C695.467 435.877 684.902 430.596 675.336 425.218L675.313 425.205L675.29 425.192C663.996 418.8 655.095 412.873 646.842 406.013C637.097 398.013 628.488 388.386 621.175 377.4C614.482 367.532 608.446 356.001 603.1 343.231C598.206 331.634 594.701 320.233 592.068 311.199L592.048 311.132L592.029 311.064C591.475 309.111 590.955 307.224 590.465 305.403C588.841 306.758 587.227 308.172 585.634 309.638L585.603 309.667L585.572 309.695C577.004 317.508 569.047 326.872 561.974 337.611C553.113 351.108 546.516 366.083 542.392 382.137L542.387 382.156L542.382 382.176C537.347 401.605 536.498 422.549 539.953 442.559L539.965 442.625L539.976 442.692C543.633 464.8 553.092 488.084 566.642 507.859L566.665 507.893L566.688 507.926C579.735 527.145 597.771 545.309 618.989 560.281C630.441 568.235 640.635 574.341 649.964 578.961L650.188 579.072L650.409 579.189C651.717 579.881 654.387 581.218 656.954 582.491C658.193 583.104 659.34 583.67 660.188 584.088L661.519 584.743L661.543 584.755L664.676 586.25L664.723 586.272C667.273 587.473 669.147 588.418 670.882 589.293C672.527 590.115 673.347 590.524 674.164 590.883L674.886 591.183L675.505 591.495C691.841 599.735 703.747 606.039 713.93 611.895C725.188 618.342 735.494 624.909 745.378 631.949L745.388 631.956L745.397 631.963C756.136 639.632 764.864 647.066 772.488 655.263L772.505 655.281L772.522 655.3C781.212 664.694 788.636 675.486 794.709 687.341C800.158 697.866 804.717 709.818 808.728 723.536L808.752 723.617L808.775 723.7C812.358 736.375 814.577 748.198 816.039 756.546L816.068 756.712L816.094 756.878C816.233 757.761 816.368 758.635 816.501 759.499C816.56 759.89 816.62 760.278 816.678 760.665C817.831 759.796 818.965 758.907 820.086 757.995L820.175 757.923L820.265 757.851C829.443 750.564 838.032 741.811 845.725 731.846L845.76 731.8L845.796 731.754C855.678 719.114 863.533 704.778 869.124 689.121C875.851 670.282 878.861 649.584 877.786 629.407L877.785 629.385L877.784 629.362C876.641 606.932 870.313 583.653 859.339 562.242L859.329 562.223Z"
            fill="url(#paint1_linear_4367_34258)"
          />
          <defs>
            <linearGradient
              id="paint0_linear_4367_34258"
              x1="748"
              y1="-1.18138e-06"
              x2="403"
              y2="1024"
              gradientUnits="userSpaceOnUse"
            >
              <stop stopColor="white" />
              <stop offset="1" stopColor="#FAF2F9" />
            </linearGradient>
            <linearGradient
              id="paint1_linear_4367_34258"
              x1="877.841"
              y1="250.557"
              x2="280.733"
              y2="896.783"
              gradientUnits="userSpaceOnUse"
            >
              <stop stop-color="#47E5BC" />
              <stop offset="1" stop-color="#5C95FF" />
            </linearGradient>
          </defs>
        </svg>
        <p className="ml-4 text-lg font-bold">Animus</p>
      </a>
      <div className="w-full px-2">
        {/* BASIC */}
        <div className="flex flex-col items-center w-full mt-3 border-t border-gray-700">
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Dashboard</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded bg-gray-700 hover:bg-gray-700 text-gray-200 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Storage</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Gateways</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.39-2.823 1.07-4"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">CID Tools</span>
          </a>
        </div>
        {/* ADVANCED */}
        <div className="flex flex-col items-center w-full mt-3 border-t border-gray-700">
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Automation</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M20.618 5.984A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016zM12 9v2m0 4h.01"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">API Access</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Private Networks</span>
          </a>
        </div>
        {/* DOCS SECTION */}
        <div className="flex flex-col items-center w-full mt-3 border-t border-gray-700">
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Organization</span>
          </a>
          <a
            className="flex items-center w-full h-12 px-3 mt-2 rounded hover:bg-gray-700 hover:text-gray-300"
            href="#"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
              />
            </svg>
            <span className="ml-2 text-sm font-medium">Documentation</span>
          </a>
        </div>
      </div>
      <AccountMenu />
    </aside>
  );
}

function AccountMenu() {
  return (
    <div className="flex items-center justify-center w-full h-16 mt-auto bg-gray-800 hover:bg-gray-700 hover:text-gray-50">
      <Menu as="div">
        <div>
          <Menu.Button
            aria-label="account menu"
            className="flex items-center w-full dark:text-gray-200 hover:text-gray-50 dark:hover:text-gray-400 focus:text-gray-50 dark:focus:text-gray-400 focus:outline-none"
          >
            <svg
              className="w-6 h-6 stroke-current"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span className="text-sm font-medium p-2">Account</span>
          </Menu.Button>
        </div>
        <Menu.Items className="absolute -translate-y-full mt-1 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 divide-y divide-gray-100 focus:outline-none">
          <div className="px-1 py-1 ">
            <Menu.Item>
              {({ active }) => (
                <Link href="/settings/profile">
                  <a
                    className={`${
                      active ? 'bg-blue-500 text-white' : 'text-gray-900'
                    } group flex rounded-md items-center w-full px-2 py-2 text-sm`}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="gray"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                      />
                    </svg>
                    <span className="px-2 font-semi">Settings</span>
                  </a>
                </Link>
              )}
            </Menu.Item>
          </div>
          <div className="px-1 py-1">
            <Menu.Item>
              {({ active }) => (
                <Link href="/settings/billing">
                  <a
                    className={`${
                      active ? 'bg-blue-500 text-white' : 'text-gray-900'
                    } group flex rounded-md items-center w-full px-2 py-2 text-sm`}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="green"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                      />
                    </svg>
                    <span className="px-2 font-semi">Billing</span>
                  </a>
                </Link>
              )}
            </Menu.Item>
            <Menu.Item>
              {({ active }) => (
                <Link href="/settings/keys">
                  <a
                    className={`${
                      active ? 'bg-blue-500 text-white' : 'text-gray-900'
                    } group flex rounded-md items-center w-full px-2 py-2 text-sm`}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="gray"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                      />
                    </svg>
                    <span className="px-2 font-semi">API Access</span>
                  </a>
                </Link>
              )}
            </Menu.Item>
          </div>
          <div className="px-1 py-1">
            <Menu.Item>
              {({ active }) => (
                <button
                  className={`${
                    active ? 'bg-blue-500 text-white' : 'text-gray-900'
                  } group flex rounded-md items-center w-full px-2 py-2 text-sm`}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="fuchsia"
                    strokeWidth={2}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                    />
                  </svg>
                  <span className="px-1">Sign Out</span>
                </button>
              )}
            </Menu.Item>
          </div>
        </Menu.Items>
      </Menu>
    </div>
  );
}
